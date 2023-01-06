/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package manager

import (
	"github.com/volcengine/datatester-go-sdk/meta/config"
)

type MetaManager interface {
	SetConfig(meta []byte)
	GetConfig() *config.ProductConfig
}
