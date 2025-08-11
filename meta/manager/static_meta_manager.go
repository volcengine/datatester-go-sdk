/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package manager

import (
	"github.com/volcengine/datatester-go-sdk/log"
	"github.com/volcengine/datatester-go-sdk/meta/config"
	"github.com/volcengine/datatester-go-sdk/utils"
)

type StaticMetaManager struct {
	metafile      string
	productConfig *config.ProductConfig
}

func NewStaticMetaManager(metafile []byte) *StaticMetaManager {
	s := &StaticMetaManager{
		metafile: string(metafile),
	}
	s.SetConfig(metafile)
	return s
}

func (s *StaticMetaManager) SetConfig(metafile []byte) {
	productConfig := &config.ProductConfig{}
	err := utils.NumberJsonApi.Unmarshal(metafile, &productConfig)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("product config json loads err")
		return
	}
	productConfig.Init()
	s.productConfig = productConfig
}

func (s *StaticMetaManager) GetConfig() *config.ProductConfig {
	return s.productConfig
}

var _ MetaManager = &StaticMetaManager{}
