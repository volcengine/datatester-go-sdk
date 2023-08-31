/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package event

import (
	//"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/volcengine/datatester-go-sdk/event/model"
	"github.com/volcengine/datatester-go-sdk/log"
)

var (
	instance Collector
	once     sync.Once
)

func GetInstance() Collector {
	once.Do(func() {
		instance = NewAsyncCollector(GetConfig())
	})
	return instance
}

type Collector interface {
	CollectEvents(events *model.ExposureEvents) error
	Stop()
}

type ResponseData struct {
	Status int `json:"e"`
}

type AsyncCollector struct {
	mscUrl        string
	appKey        string
	mscHttpClient *http.Client
	workerNum     int
	msgChan       chan *model.ExposureEvents
	exitChan      chan struct{}
	batchSize     int
	lingerTime    time.Duration
	successCount  uint64
	failedCount   uint64
}

func NewAsyncCollector(config Config) Collector {
	asyncConfig := config.AsyncConfig

	collector := &AsyncCollector{
		mscUrl: config.MscUrl,
		appKey: config.AppKey,
		mscHttpClient: &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, config.HttpDialTimeout)
				},
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: config.HttpMaxIdleConnPerHost,
				MaxConnsPerHost:     config.HttpMaxConnPerHost,
				IdleConnTimeout:     300 * time.Second,
				// skip ssl verify
				// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: config.HttpTotalTimeout,
		},
		workerNum:    asyncConfig.WorkerNum,
		msgChan:      make(chan *model.ExposureEvents, asyncConfig.ChannelSize),
		exitChan:     make(chan struct{}),
		batchSize:    asyncConfig.BatchSize,
		lingerTime:   asyncConfig.LingerTime,
		successCount: 0,
		failedCount:  0,
	}
	collector.startStatistics()
	collector.startWorkers()
	return collector
}

func (s *AsyncCollector) startWorkers() {
	for i := 0; i < s.workerNum; i++ {
		go func() {
			batchEvents := make([]*model.ExposureEvents, 0, s.batchSize)
			lingerTimer := time.NewTimer(s.lingerTime)
			for {
				select {
				case <-s.exitChan:
					lingerTimer.Stop()
					log.Warn("receive exit sign, worker stop")
					return
				case event := <-s.msgChan:
					batchEvents = append(batchEvents, event)
					if len(batchEvents) >= s.batchSize {
						if err := s.send(batchEvents); err != nil {
							atomic.AddUint64(&s.failedCount, uint64(len(batchEvents)))
							log.WithFields(log.Fields{"err": err}).ErrorF(
								"[%d] events dispatcher failed, events: %v", len(batchEvents), batchEvents)
						} else {
							atomic.AddUint64(&s.successCount, uint64(len(batchEvents)))
							log.InfoF("[%d] events dispatcher success", len(batchEvents))
						}
						batchEvents = make([]*model.ExposureEvents, 0, s.batchSize)
						lingerTimer.Reset(s.lingerTime)
					}
				case <-lingerTimer.C:
					if len(batchEvents) > 0 {
						if err := s.send(batchEvents); err != nil {
							atomic.AddUint64(&s.failedCount, uint64(len(batchEvents)))
							log.WithFields(log.Fields{"err": err}).ErrorF(
								"[%d] events dispatcher failed, events: %v", len(batchEvents), batchEvents)
						} else {
							atomic.AddUint64(&s.successCount, uint64(len(batchEvents)))
							log.InfoF("[%d] events dispatcher success", len(batchEvents))
						}
						batchEvents = make([]*model.ExposureEvents, 0, s.batchSize)
					}
					lingerTimer.Reset(s.lingerTime)
				}
			}
		}()
	}
}

func (s *AsyncCollector) startStatistics() {
	go func() {
		// print every 30 seconds
		printTimer := time.NewTimer(30 * time.Second)
		// reset the counter at 24 o'clock every day
		next := time.Now().Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		resetTimer := time.NewTimer(next.Sub(time.Now()))
		for {
			select {
			case <-printTimer.C:
				log.InfoF(
					"[Statistics] current success count [%d], current failed count [%d], current lag [%d]",
					atomic.LoadUint64(&s.successCount), atomic.LoadUint64(&s.failedCount), len(s.msgChan))
				printTimer.Reset(30 * time.Second)
			case <-resetTimer.C:
				log.InfoF("[Statistics] total success count [%d], total failed count [%d]",
					atomic.SwapUint64(&s.successCount, 0), atomic.SwapUint64(&s.failedCount, 0))
				next = time.Now().Add(time.Hour * 24)
				next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
				resetTimer.Reset(next.Sub(time.Now()))
			}
		}
	}()
}

func (s *AsyncCollector) CollectEvents(events *model.ExposureEvents) error {
	s.msgChan <- events
	return nil
}

func (s *AsyncCollector) Stop() {
	close(s.exitChan)
}

func (s *AsyncCollector) send(events []*model.ExposureEvents) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.mscUrl, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("User-Agent", "GoSDK")
	if s.appKey != "" {
		req.Header.Set("X-MCS-AppKey", s.appKey)
	}
	if resp, err := s.mscHttpClient.Do(req); err == nil {
		defer func() {
			_ = resp.Body.Close()
		}()
		if resp.StatusCode != 200 {
			return fmt.Errorf("response code: %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %v", err)
		}
		var responseData ResponseData
		if err = json.Unmarshal(body, &responseData); err != nil {
			return fmt.Errorf("read response body: %v", err)
		}
		if responseData.Status != 0 {
			return fmt.Errorf("collected response: %s", string(body))
		}
	} else {
		return err
	}
	return nil
}
