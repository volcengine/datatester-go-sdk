/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllowWithoutFilter(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("allow_without_filter", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	result, _ = abClient.Activate("allow_without_filter", "test_user",
		"trackId", nil, attributes)
	assert.True(t, result.(bool) == true)
	result, _ = abClient.Activate("allow_without_filter", "test_user_02",
		"trackId", nil, attributes)
	assert.True(t, result.(bool) == false)
}

func TestAllowWithFilter(t *testing.T) {
	abClient, attributes := setUp()
	result, _ := abClient.Activate("allow_with_filter", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = nil
	result, _ = abClient.Activate("allow_with_filter", "test_user",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	result, _ = abClient.Activate("allow_with_filter", "test_user_02",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	attributes["str_param"] = ""
	result, _ = abClient.Activate("allow_with_filter", "test_user",
		"trackId", nil, attributes)
	assert.True(t, result.(float64) == 0)
	attributes["str_param"] = "str1"
	result, _ = abClient.Activate("allow_with_filter", "test_user_02",
		"trackId", nil, attributes)
	assert.True(t, result.(float64) == 1)
}
