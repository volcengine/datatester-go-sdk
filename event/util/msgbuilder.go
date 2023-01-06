/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package util

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/volcengine/datatester-go-sdk/consts"
	"github.com/volcengine/datatester-go-sdk/event/model"
	"strconv"
	"time"
)

var (
	paramsFormat = "{\"datatester_sdk_version\":\"1.0.3\",\"datatester_sdk_language\":\"golang\"," +
		"\"$user_unique_id_type\":\"%s\"}"
)

func CreateHeader(appId uint32, tz int32) *model.Header {
	return &model.Header{
		AppId:    proto.Uint32(appId),
		Timezone: proto.Int32(tz),
	}
}

func CreateExposureEvent(vid, uuidType string) *model.Event {
	exposureEvent := &model.Event{
		Event:        proto.String("abtest_exposure"),
		LocalTimeMs:  proto.Uint64(uint64(time.Now().UnixNano() / 1e6)),
		AbSdkVersion: proto.String(vid),
		Params:       proto.String(fmt.Sprintf(paramsFormat, uuidType)),
	}
	return exposureEvent
}

func CreateUser(trackId, uuidType string, attributes map[string]interface{},
	enableAnonymousUser, isSaasEnv bool) *model.User {
	if enableAnonymousUser && trackId == "" && attributes != nil {
		return CreateAnonymousUser(trackId, uuidType, attributes, isSaasEnv)
	}
	return &model.User{
		UserUniqueId: proto.String(trackId),
		UuidType:     proto.String(uuidType),
	}
}

func CreateAnonymousUser(trackId, uuidType string, attributes map[string]interface{}, isSaasEnv bool) *model.User {
	user := &model.User{
		UserUniqueId: proto.String(trackId),
		UuidType:     proto.String(uuidType),
	}
	deviceId := getIdByType(consts.DeviceId, attributes)
	if isSaasEnv {
		if deviceId != -1 {
			user.DeviceId = proto.Int64(deviceId)
		}
		bdDid := getBdDid(attributes)
		if bdDid != "" {
			user.BdDid = proto.String(bdDid)
		}
	} else {
		if deviceId != -1 {
			user.DeviceId = proto.Int64(deviceId)
		} else {
			bdDid := getBdDid(attributes)
			if bdDid != "" {
				int64Val, err := strconv.ParseInt(bdDid, 10, 64)
				if err == nil {
					user.DeviceId = proto.Int64(int64Val)
				}
			}
		}
	}
	webId := getIdByType(consts.WebId, attributes)
	if webId != -1 {
		user.WebId = proto.Int64(webId)
	}
	return user
}

func getIdByType(idType string, attributes map[string]interface{}) int64 {
	val, exist := attributes[idType]
	if !exist || val == nil {
		return -1
	}
	switch val.(type) {
	case int:
		return int64(val.(int))
	case int32:
		return int64(val.(int32))
	case int64:
		return val.(int64)
	case float64:
		return int64(val.(float64))
	case string:
		int64Val, err := strconv.ParseInt(val.(string), 10, 64)
		if err == nil {
			return int64Val
		}
		return -1
	default:
		return -1
	}
}

func getBdDid(attributes map[string]interface{}) string {
	val, exist := attributes[consts.BdDid]
	if !exist || val == nil {
		return ""
	}
	if strVal, ok := val.(string); ok {
		return strVal
	}
	return ""
}

func GetTimezone() int32 {
	utc, _ := strconv.Atoi(time.Now().UTC().Format("2006-01-02 15:04:05")[11:13])
	cur, _ := strconv.Atoi(time.Now().Format("2006-01-02 15:04:05")[11:13])
	ans := cur - utc
	if ans <= -12 {
		ans = ans + 24
	}
	return int32(ans)
}
