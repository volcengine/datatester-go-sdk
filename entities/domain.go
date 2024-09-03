package entities

import (
	"strings"

	"github.com/volcengine/datatester-go-sdk/distributor/bucketer"
)

type Domain struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	HashStrategy string `json:"hash_strategy"`
	Begin        uint32 `json:"begin"`
	Length       uint32 `json:"length"`
}

// Hit 方法检查给定的 decisionID 是否在当前 Domain 的索引范围内
func (d Domain) Hit(decisionID string) bool {
	// 通过连接 decisionID 和 d.Name 获取流量桶索引
	index, err := bucketer.NewMmh3BucketService().GetTrafficBucketIndex(
		strings.Join([]string{decisionID, d.Name}, ":"))
	// 如果获取索引过程中出现错误，返回 false
	if err != nil {
		return false
	}
	// 如果索引值在 d.Begin 和 d.Begin + d.Length 之间，返回 true，否则返回 false
	return index >= d.Begin && index < d.Begin+d.Length
}
