/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumberFilterParam(t *testing.T) {
	abClient, attributes := setUp()
	attributes["str_param"] = "fgh"
	result, _ := abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 56789
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 123.45
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
	attributes["number_param"] = 99999
	result, _ = abClient.Activate("filter_param", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "b")
}

func TestNumberFilterParam02(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 12345
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 12345.001
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["number_param"] = 6789
	result, _ = abClient.Activate("filter_param_02", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
}

func TestNumberFilterParam03(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 345.56
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["number_param"] = 346
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["number_param"] = 123.45
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["number_param"] = 120
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["number_param"] = 180
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["number_param"] = 300
	result, _ = abClient.Activate("filter_param_03", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
}
