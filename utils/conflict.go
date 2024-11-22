package utils

import (
	"strconv"
	"strings"
)

func IsHigherPriorityConfig(existConf, newConf map[string]interface{}) bool {
	vid1 := existConf["vid"].(string)
	vid2 := newConf["vid"].(string)
	v1, err1 := strconv.ParseInt(vid1, 10, 64)
	v2, err2 := strconv.ParseInt(vid2, 10, 64)
	// 转换失败，则按照字符串对比，一般不会出现
	if err1 != nil || err2 != nil {
		return strings.Compare(vid2, vid1) < 0
	}
	return v2 < v1
}
