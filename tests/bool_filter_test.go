/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolFilterParam(t *testing.T) {
	abClient, attributes := setUp()
	attributes["str_param"] = "fgh"
	result, _ := abClient.Activate("filter_param", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result == nil)
	attributes["bool_param"] = true
	result, _ = abClient.Activate("filter_param", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "b")
	attributes["bool_param"] = false
	result, _ = abClient.Activate("filter_param", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "b")
}

func TestBoolFilterParam02(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_02", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result == nil)
	attributes["bool_param"] = true
	result, _ = abClient.Activate("filter_param_02", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["bool_param"] = false
	result, _ = abClient.Activate("filter_param_02", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "a")
}

func TestBoolFilterParam03(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("filter_param_03", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result == nil)
	attributes["bool_param"] = true
	result, _ = abClient.Activate("filter_param_03", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "a")
	attributes["bool_param"] = false
	result, _ = abClient.Activate("filter_param_03", "decisionId", "trackId",
		nil, attributes)
	assert.True(t, result.(string) == "a")
}
