/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package tests

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAssociationExperiments(t *testing.T) {
	abClient, attributes := setUp()
	i := 0
	j := 0
	for count := 0; i < 3 || j < 3; count++ {
		decisionId := "decisionId" + strconv.Itoa(count)
		childMap, _ := abClient.GetExperimentConfigs("77779", decisionId, attributes)
		attributes = map[string]interface{}{}
		fatherMap, _ := abClient.GetExperimentConfigs("77775", decisionId, attributes)
		attributes = map[string]interface{}{}
		if childMap == nil {
			assert.True(t, fatherMap == nil)
			i++
		} else {
			assert.True(t, fatherMap != nil)
			j++
		}
	}
	i = 0
	j = 0
	for count := 0; i < 3 || j < 3; count++ {
		decisionId := "decisionId" + strconv.Itoa(count)
		childMap, _ := abClient.GetExperimentConfigs("77780", decisionId, attributes)
		attributes = map[string]interface{}{}
		fatherMap, _ := abClient.GetExperimentConfigs("77775", decisionId, attributes)
		attributes = map[string]interface{}{}
		if childMap == nil {
			assert.True(t, fatherMap != nil)
			i++
		} else {
			assert.True(t, fatherMap == nil)
			j++
		}
	}
}
