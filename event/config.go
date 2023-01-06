/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package event

import (
	"time"
)

const (
	defaultTraceUrl               = "https://mcs.ctobsnssdk.com/v2/event/list"
	defaultTracePath              = "/v2/event/list"
	defaultHttpDialTimeout        = 1 * time.Second
	defaultHttpTotalTimeout       = 5 * time.Second
	defaultHttpMaxIdleConnPerHost = 10
	defaultWorkerNum              = 20
	defaultChanSize               = 20000
	defaultBatchSize              = 50
	defaultLingerTime             = 1 * time.Second
)

type Config struct {
	MscUrl                  string
	AppKey                  string
	HttpDialTimeout         time.Duration
	HttpTotalTimeout        time.Duration
	HttpMaxIdleConnsPerHost int
	AsyncConfig             *AsyncConfig
}

type AsyncConfig struct {
	ChannelSize int
	WorkerNum   int
	BatchSize   int
	LingerTime  time.Duration
}

type AnonymousUserConfig struct {
	Enable bool
	IsSaas bool
}

func NewDefaultConfig() Config {
	asyncConfig := &AsyncConfig{
		ChannelSize: defaultChanSize,
		WorkerNum:   defaultWorkerNum,
		BatchSize:   defaultBatchSize,
		LingerTime:  defaultLingerTime,
	}
	return Config{
		MscUrl:                  defaultTraceUrl,
		HttpDialTimeout:         defaultHttpDialTimeout,
		HttpTotalTimeout:        defaultHttpTotalTimeout,
		HttpMaxIdleConnsPerHost: defaultHttpMaxIdleConnPerHost,
		AsyncConfig:             asyncConfig,
	}
}

// params support change
var (
	traceUrl            = defaultTraceUrl
	gWorkerNum          = defaultWorkerNum
	anonymousUserConfig = &AnonymousUserConfig{
		Enable: false,
		IsSaas: true,
	}
)

func SetTraceUrl(traceHost string) {
	traceUrl = traceHost + defaultTracePath
}

func SetWorkerNum(workerNum int) {
	if workerNum > 0 {
		gWorkerNum = workerNum
	}
}

func SetAnonymousUserConfig(enable bool, isSaas bool) {
	anonymousUserConfig.Enable = enable
	anonymousUserConfig.IsSaas = isSaas
}

func GetConfig() Config {
	config := NewDefaultConfig()
	config.MscUrl = traceUrl
	config.AsyncConfig.WorkerNum = gWorkerNum
	return config
}

func GetAnonymousUserConfig() *AnonymousUserConfig {
	return anonymousUserConfig
}
