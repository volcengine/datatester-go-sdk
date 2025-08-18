package utils

import (
	"strconv"
	"strings"

	et "github.com/volcengine/datatester-go-sdk/entities"
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

// IsHigherPriorityConfig2 判断new的优先级是否高于exist
func IsHigherPriorityConfig2(existConf, newConf map[string]interface{},
	existExpID string, newExpID string, experimentMap map[string]*et.Experiment) bool {
	existExp, exist1 := experimentMap[existExpID]
	newExp, exist2 := experimentMap[newExpID]
	if !exist1 {
		return true
	}
	if !exist2 {
		return false
	}
	if !existExp.IsCodingCampaign() && !newExp.IsCodingCampaign() {
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
	// 个性化实验优先级高于普通实验
	if existExp.IsCodingCampaign() && !newExp.IsCodingCampaign() {
		return false
	}
	// 个性化实验优先级高于普通实验
	if !existExp.IsCodingCampaign() && newExp.IsCodingCampaign() {
		return true
	}
	// 都是个性化实验，如果fvid不同，则按照fvid比较，fvid大的优先级高
	// 如果fvid相同，则根据子实验优先级进行对比
	existFVid1, e1 := strconv.Atoi(existConf["f_vid"].(string))
	if e1 != nil {
		existFVid1 = -1
	}
	newFVid2, e2 := strconv.Atoi(newConf["f_vid"].(string))
	if e2 != nil {
		newFVid2 = -1
	}
	if newFVid2 != existFVid1 {
		return newFVid2 > existFVid1
	}
	return newExp.FlightPriority > existExp.FlightPriority
}
