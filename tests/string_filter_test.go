/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFilterParam(t *testing.T) {
	abClient, attributes := setUp()
	attributes["str_param"] = "fgh"
	result, _ := abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes = map[string]interface{}{}
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
	attributes["str_param"] = "jll"
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
	attributes["str_param"] = "aaa"
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
}

func TestStringFilterParam02(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = "5.6.7"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["str_param"] = "6.6.7"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["str_param"] = "2.3.4"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["str_param"] = "2.2.3"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["str_param"] = "2.3.5"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = "5.6.6"
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
}

func TestStringFilterParam03(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = "str1"
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = "str2"
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = ""
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["str_param"] = "str3"
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
}

func TestStringFilterParam04(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_04", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = ""
	result, _ = abClient.Activate("filter_param_04", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
	attributes["str_param"] = "str1"
	result, _ = abClient.Activate("filter_param_04", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
}
