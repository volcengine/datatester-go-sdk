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

func TestFatherChildExperiments(t *testing.T) {
	abClient, attributes := setUp()
	i := 0
	j := 0
	for count := 0; i < 3 || j < 3; count++ {
		decisionId := "decisionId" + strconv.Itoa(count)
		childMap, _ := abClient.GetExperimentConfigs("77774", decisionId, attributes)
		if childMap == nil {
			continue
		}
		if childMap["child_01"]["vid"] == "120311" {
			fatherMap, _ := abClient.GetExperimentConfigs("77773", decisionId, attributes)
			assert.True(t, fatherMap["father"]["vid"] == "120309")
			i++
		} else {
			fatherMap, _ := abClient.GetExperimentConfigs("77773", decisionId, attributes)
			assert.True(t, fatherMap["father"]["vid"] == "120310")
			j++
		}
	}
}
