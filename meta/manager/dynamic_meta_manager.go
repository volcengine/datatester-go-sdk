/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package manager

import (
	"errors"
	"fmt"
	"github.com/volcengine/datatester-go-sdk/event"
	"github.com/volcengine/datatester-go-sdk/log"
	"github.com/volcengine/datatester-go-sdk/utils/httper"
	"strconv"
	"strings"
	"time"
)

// meta url
const (
	DefaultMetaUrl  = "https://datarangers.com.cn/abmeta/v2/get_abtest_info?token=%s&sideType=%s"
	DefaultMetaPath = "/abmeta/v2/get_abtest_info?token=%s&sideType=%s"
)

const (
	DefaultUpdateInterval = 1 * time.Minute
	LastModified          = "Modify_time"
)

type MetaSideType string

const (
	SERVER MetaSideType = "server"
)

func AdaptHost(host string) string {
	return strings.TrimRight(host, "/")
}

type DynamicMetaManager struct {
	requester     *httper.HRequester
	token         string
	fetchInterval time.Duration
	baseUrl       string
	fetchUrl      string
	LastModified  int
	StaticMetaManager
}

type MetaOptionFunc func(*DynamicMetaManager)

func WithRequester(requester *httper.HRequester) MetaOptionFunc {
	return func(d *DynamicMetaManager) {
		d.requester = requester
	}
}

func WithMetaHost(host string) MetaOptionFunc {
	return func(d *DynamicMetaManager) {
		url := AdaptHost(host) + DefaultMetaPath
		d.fetchUrl = fmt.Sprintf(url, d.token, SERVER)
	}
}

func WithTrackHost(host string) MetaOptionFunc {
	return func(d *DynamicMetaManager) {
		event.SetTraceUrl(AdaptHost(host))
	}
}

func WithFetchInterval(interval time.Duration) MetaOptionFunc {
	if interval <= 0 {
		interval = DefaultUpdateInterval
	}
	return func(d *DynamicMetaManager) {
		d.fetchInterval = interval
	}
}

func NewDynamicMetaManager(token string, options ...MetaOptionFunc) *DynamicMetaManager {
	d := DynamicMetaManager{
		requester:     httper.NewHTTPRequester(),
		token:         token,
		fetchInterval: DefaultUpdateInterval,
		baseUrl:       DefaultMetaUrl,
		fetchUrl:      fmt.Sprintf(DefaultMetaUrl, token, SERVER),
	}
	for _, opt := range options {
		opt(&d)
	}
	if err := d.fetchMeta(); err != nil {
		log.WithFields(log.Fields{"url": d.fetchUrl, "err": err}).Error("fetch meta err")
		return &d
	}
	d.startFetcher()
	return &d
}

func NewDynamicMetaManager4Test(meta string) *DynamicMetaManager {
	d := DynamicMetaManager{
		requester:     httper.NewHTTPRequester(),
		fetchInterval: DefaultUpdateInterval,
		baseUrl:       DefaultMetaUrl,
	}
	d.StaticMetaManager = *NewStaticMetaManager([]byte(meta))
	return &d
}

func (d *DynamicMetaManager) startFetcher() {
	log.InfoF("start fetch meta..., interval is %v", d.fetchInterval)
	t := time.NewTicker(d.fetchInterval)
	go func() {
		for {
			<-t.C
			_ = d.fetchMeta()
		}
	}()
}

func (d *DynamicMetaManager) fetchMeta() (err error) {
	defer func() {
		if errTemp := recover(); errTemp != nil {
			log.ErrorF("fetch meta panic: %v", errTemp)
			err = errors.New(errTemp.(string))
		}
	}()
	var response *httper.HResponse
	response, err = d.requester.Get(d.fetchUrl)
	if err != nil {
		log.WithFields(log.Fields{"url": d.fetchUrl, "err": err}).Error("fetch meta err")
		return err
	}
	log.WithFields(log.Fields{"url": d.fetchUrl, "code": response.Code}).Info(
		"fetch meta success")
	if lastModified, _ := strconv.Atoi(response.Headers.Get(LastModified)); lastModified != d.LastModified {
		d.setLastModified(lastModified)
		d.SetConfig(response.Body)
	}
	return nil
}

func (d *DynamicMetaManager) GetFetchUrl() string {
	return d.fetchUrl
}

func (d *DynamicMetaManager) SetFetchUrl(url string) {
	d.fetchUrl = url
}

func (d *DynamicMetaManager) setLastModified(lastModified int) {
	d.LastModified = lastModified
}

var _ MetaManager = &DynamicMetaManager{}
