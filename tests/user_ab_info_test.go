/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/volcengine/datatester-go-sdk/client"
	"github.com/volcengine/datatester-go-sdk/tests/mock"
	"io/ioutil"
	"testing"
)

var handler = NewTextUserAbInfoHandler()

func setUpWithHandler() (*client.AbClient, map[string]interface{}) {
	return client.NewClient4Test(mock.MetaMock, handler), make(map[string]interface{})
}

func TestFreezeVersion(t *testing.T) {
	abClient, attributes := setUpWithHandler()
	clearCache()
	resultMap, _ := abClient.GetExperimentConfigs("77773", "decisionId", attributes)
	assert.True(t, resultMap["father"]["vid"] == "120310")
	clearCache()
	handler.CreateOrUpdate("decisionId", "{\"77773\":\"120309\"}")
	resultMap, _ = abClient.GetExperimentConfigs("77773", "decisionId", attributes)
	assert.True(t, resultMap["father"]["vid"] == "120309")
	clearCache()
}

func TestFreezeExperiment(t *testing.T) {
	abClient, attributes := setUpWithHandler()
	clearCache()
	resultMap, _ := abClient.GetExperimentConfigs("77772", "decisionId", attributes)
	assert.True(t, len(resultMap) == 0)
	handler.CreateOrUpdate("decisionId", "{\"77772\":\"120307\"}")
	resultMap, _ = abClient.GetExperimentConfigs("77772", "decisionId", attributes)
	assert.True(t, resultMap["freeze"]["vid"] == "120307")
	clearCache()
	handler.CreateOrUpdate("decisionId", "{\"77772\":\"120308\"}")
	resultMap, _ = abClient.GetExperimentConfigs("77772", "decisionId", attributes)
	assert.True(t, resultMap["freeze"]["vid"] == "120308")
	clearCache()
	handler.CreateOrUpdate("decisionId", "{\"77772\":\"120309\"}")
	resultMap, _ = abClient.GetExperimentConfigs("77772", "decisionId", attributes)
	assert.True(t, len(resultMap) == 0)
	clearCache()
}

func TestActivateAndGetAllExperiment(t *testing.T) {
	abClient, attributes := setUpWithHandler()
	clearCache()
	abClient.GetAllExperimentConfigs("decisionId", attributes)
	assert.True(t, len(transferStr2Map(handler.Query("decisionId"))) == 5)
	clearCache()
	abClient.Activate("asso", "decisionId", "trackId", nil, attributes)
	assert.True(t, len(transferStr2Map(handler.Query("decisionId"))) == 1)
	clearCache()
	abClient.Activate("not_exist_key", "decisionId", "trackId", nil, attributes)
	assert.True(t, len(transferStr2Map(handler.Query("decisionId"))) == 0)
	clearCache()
	val, _ := abClient.Activate(
		"freeze", "decisionId", "trackId", nil, attributes)
	assert.True(t, val == nil)
	assert.True(t, len(transferStr2Map(handler.Query("decisionId"))) == 0)
	handler.CreateOrUpdate("decisionId", "{\"77772\":\"120307\"}")
	val, _ = abClient.Activate("freeze", "decisionId", "trackId", nil, attributes)
	assert.True(t, val.(map[string]interface{})["freeze"].(float64) == 0)
	assert.True(t, len(transferStr2Map(handler.Query("decisionId"))) == 1)
	clearCache()
}

type TextAbInfoHandler struct{}

func (u *TextAbInfoHandler) Query(decisionID string) string {
	filepath := "./mock/user_info.txt"
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	mapInfo := make(map[string]string)
	json.Unmarshal(content, &mapInfo)
	if value, ok := mapInfo[decisionID]; ok {
		return value
	}
	return ""
}

func (u *TextAbInfoHandler) CreateOrUpdate(decisionID, experiment2Variant string) bool {
	filepath := "./mock/user_info.txt"
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	mapInfo := make(map[string]string)
	err = json.Unmarshal(content, &mapInfo)
	if err != nil {
		mapInfo = make(map[string]string)
	}
	mapInfo[decisionID] = experiment2Variant
	data, _ := json.Marshal(mapInfo)
	err = ioutil.WriteFile("./mock/user_info.txt", data, 0644)
	if err != nil {
		panic(err)
	}
	return true
}

func (u *TextAbInfoHandler) NeedPersistData() bool {
	return true
}

func NewTextUserAbInfoHandler() *TextAbInfoHandler {
	return &TextAbInfoHandler{}
}

func clearCache() {
	mapInfo := make(map[string]string)
	data, _ := json.Marshal(mapInfo)
	err := ioutil.WriteFile("./mock/user_info.txt", data, 0644)
	if err != nil {
		panic(err)
	}
}

func transferStr2Map(str string) map[string]string {
	result := make(map[string]string)
	json.Unmarshal([]byte(str), &result)
	return result
}
