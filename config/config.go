/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package config

import (
	"github.com/volcengine/datatester-go-sdk/event"
	"github.com/volcengine/datatester-go-sdk/log"
	"github.com/volcengine/datatester-go-sdk/meta/manager"
	"github.com/volcengine/datatester-go-sdk/utils/httper"
	"time"
)

const (
	MetaHostCN string = "https://datarangers.com.cn"
	MetaHostSG string = "https://datarangers.com"
)

const (
	TrackHostCN string = "https://mcs.ctobsnssdk.com"
	TrackHostSG string = "https://mcs.tobsnssdk.com"
	TrackHostVA string = "https://mcs.itobsnssdk.com"
)

type Func func() (manager.MetaOptionFunc, bool)

func WithRequester(requester *httper.HRequester) Func {
	return func() (manager.MetaOptionFunc, bool) {
		return manager.WithRequester(requester), true
	}
}

func WithMetaHost(host string) Func {
	return func() (manager.MetaOptionFunc, bool) {
		return manager.WithMetaHost(host), true
	}
}

func WithTrackHost(host string) Func {
	return func() (manager.MetaOptionFunc, bool) {
		return manager.WithTrackHost(host), true
	}
}

func WithFetchInterval(interval time.Duration) Func {
	return func() (manager.MetaOptionFunc, bool) {
		return manager.WithFetchInterval(interval), true
	}
}

func WithWorkerNumOnce(workerNum int) Func {
	return func() (manager.MetaOptionFunc, bool) {
		event.SetWorkerNum(workerNum)
		return nil, false
	}
}

func WithAnonymousConfig(enable bool, isSaasEnv bool) Func {
	return func() (manager.MetaOptionFunc, bool) {
		event.SetAnonymousUserConfig(enable, isSaasEnv)
		return nil, false
	}
}

func WithLogger(l log.Logger) Func {
	return func() (manager.MetaOptionFunc, bool) {
		log.InitGlobalLogger(l)
		return nil, false
	}
}
