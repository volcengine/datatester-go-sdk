/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package bucketer

import (
	"github.com/twmb/murmur3"
	"github.com/volcengine/datatester-go-sdk/log"
)

const (
	DefaultHashSeed = 0
	MaxTrafficValue = 1000
)

type Mmh3BucketService struct {
	hashSeed uint32
}

func NewMmh3BucketService() *Mmh3BucketService {
	return &Mmh3BucketService{
		hashSeed: DefaultHashSeed,
	}
}

func (m *Mmh3BucketService) generateHashCode(hashKey string) (int32, error) {
	mm3 := murmur3.SeedNew32(m.hashSeed)
	if _, err := mm3.Write([]byte(hashKey)); err != nil {
		log.WithFields(log.Fields{
			"hashKey": hashKey,
			"err":     err,
		}).Error("generate hashcode for hash key")
		return 0, err
	}
	hashCode := int32(mm3.Sum32())
	return hashCode, nil
}

func (m *Mmh3BucketService) generateBucketIndex(hashCode int32) uint32 {
	index := hashCode % MaxTrafficValue
	if index < 0 {
		index += MaxTrafficValue
	}
	return uint32(index)
}

func (m *Mmh3BucketService) GetTrafficBucketIndex(hashKey string) (uint32, error) {
	hashCode, err := m.generateHashCode(hashKey)
	if err != nil {
		return 0, err
	}
	index := m.generateBucketIndex(hashCode)
	return index, nil
}

var _ BucketService = &Mmh3BucketService{}
