/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package bucketer

type BucketService interface {
	generateHashCode(hashKey string) (int32, error)
	generateBucketIndex(hashCode int32) uint32
}
