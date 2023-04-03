/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/volcengine/datatester-go-sdk/client"
	"github.com/volcengine/datatester-go-sdk/config"
	"github.com/volcengine/datatester-go-sdk/event/model"
	"github.com/volcengine/datatester-go-sdk/event/util"
	"testing"
)

func TestAnonymousUserConfig(t *testing.T) {
	abClient := client.NewClient("3ac7ff2fad9d005f97857596a9078aff", config.WithMetaHost(config.MetaHostCN),
		config.WithTrackHost(config.TrackHostCN), config.WithAnonymousConfig(false, true))
	assert.True(t, abClient != nil)
	val, _ := abClient.Activate(
		"go", "decisionId", "trackId", nil, map[string]interface{}{})
	if val != nil {
		assert.True(t, val.(string) == "v1" || val.(string) == "v2")
	}
}

func TestExposureEventString(t *testing.T) {
	e := util.CreateExposureEvent("123", "test_type")
	mockTimeMs := uint64(1672893222638)
	e.LocalTimeMs = &mockTimeMs
	user := util.CreateUser("trackId", "test_type", map[string]interface{}{},
		false, true)
	events := &model.ExposureEvents{
		User:   user,
		Header: util.CreateHeader(uint32(123456), util.GetTimezone()),
		Events: []*model.Event{e},
	}
	str := fmt.Sprintf("events = [%v]", events)
	assert.True(t, str == "events = [{user: {user_unique_id: ***, user_unique_id_type: test_type}, header: {app_id: 123456, timezone: 8}, events: [{event: abtest_exposure, params: {\"datatester_sdk_version\":\"1.0.4\",\"datatester_sdk_language\":\"golang\",\"$user_unique_id_type\":\"test_type\"}, local_time_ms: 1672893222638, ab_sdk_version: 123}]}]")
}

func TestCreateUser(t *testing.T) {
	trackId := "trackId"
	user := util.CreateUser(trackId, "", map[string]interface{}{},
		false, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.UuidType == "")
	assert.True(t, user.DeviceId == nil)
	assert.True(t, user.BdDid == nil)
	assert.True(t, user.WebId == nil)
}

func TestSaasCreateAnonymousUser(t *testing.T) {
	trackId := "trackId"
	attributes := make(map[string]interface{})
	attributes["device_id"] = 123
	user := util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	trackId = ""
	attributes["device_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	attributes["device_id"] = int64(1234)
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 1234)
	attributes["device_id"] = int64(9223372036854775807)
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 9223372036854775807)
	attributes = make(map[string]interface{})
	attributes["web_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.WebId == nil)
	attributes["web_id"] = int64(1234)
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.WebId == 1234)
	attributes["web_id"] = int64(9223372036854775807)
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.WebId == 9223372036854775807)
	attributes = make(map[string]interface{})
	attributes["bddid"] = nil
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = 1234
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = ""
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = "trackId"
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.BdDid == "trackId")
	attributes = make(map[string]interface{})
	attributes["bddid"] = "trackId"
	attributes["web_id"] = 123456
	attributes["device_id"] = 67891011
	user = util.CreateUser(trackId, "", attributes, true, true)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 67891011)
	assert.True(t, *user.BdDid == "trackId")
	assert.True(t, *user.WebId == 123456)
}

func TestNotSaasCreateAnonymousUser(t *testing.T) {
	trackId := "trackId"
	attributes := make(map[string]interface{})
	attributes["device_id"] = 123
	user := util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	trackId = ""
	attributes["device_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	attributes["device_id"] = 1234
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 1234)
	attributes["device_id"] = 9223372036854775807
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 9223372036854775807)
	attributes = make(map[string]interface{})
	attributes["web_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.WebId == nil)
	attributes["web_id"] = 1234
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.WebId == 1234)
	attributes["web_id"] = 9223372036854775807
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.WebId == 9223372036854775807)
	attributes = make(map[string]interface{})
	attributes["bddid"] = nil
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = 1234
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = ""
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = "trackId"
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	assert.True(t, user.BdDid == nil)
	attributes["bddid"] = "9223372036854775807"
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 9223372036854775807)
	assert.True(t, user.BdDid == nil)
	attributes = make(map[string]interface{})
	attributes["bddid"] = "12345"
	attributes["web_id"] = 56789
	attributes["device_id"] = 101112
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 101112)
	assert.True(t, *user.WebId == 56789)
	assert.True(t, user.BdDid == nil)
	attributes = make(map[string]interface{})
	attributes["bddid"] = "12345"
	attributes["web_id"] = 56789
	attributes["device_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, *user.DeviceId == 12345)
	assert.True(t, *user.WebId == 56789)
	assert.True(t, user.BdDid == nil)
	attributes = make(map[string]interface{})
	attributes["bddid"] = ""
	attributes["web_id"] = 56789
	attributes["device_id"] = nil
	user = util.CreateUser(trackId, "", attributes, true, false)
	assert.True(t, *user.UserUniqueId == trackId)
	assert.True(t, user.DeviceId == nil)
	assert.True(t, *user.WebId == 56789)
	assert.True(t, user.BdDid == nil)
}
