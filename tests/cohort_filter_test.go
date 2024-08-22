package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeCohortClient struct {
	fakeCohortIds []string
}

func (c *fakeCohortClient) UseCohort(decisionId string, cohortIds []string) []string {
	return c.fakeCohortIds
}

func TestExpCohortFilter(t *testing.T) {
	abClient, attributes := setUp()

	// no cohort handler
	result, _ := abClient.Activate("cohort_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)

	// cohort handler, empty res
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{}})
	result, _ = abClient.Activate("cohort_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)

	// cohort handler, match filter
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123"}})
	result, _ = abClient.Activate("cohort_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	configMap, _ := abClient.GetExperimentConfigs("77790", "decisionId", attributes)
	assert.NotNil(t, configMap["cohort_config"])

	// cohort handler, not match filter
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123", "456"}})
	result, _ = abClient.Activate("cohort_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result == nil)
	configMap, _ = abClient.GetExperimentConfigs("77790", "decisionId", attributes)
	assert.Nil(t, configMap["cohort_config"])
}

func TestFeatureCohortFilter(t *testing.T) {
	abClient, attributes := setUp()

	// no cohort handler
	configMap, _ := abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.Nil(t, configMap["cohort_feature"])
	configMap, _ = abClient.GetFeatureConfigs("10301", "decisionId", attributes)
	assert.Nil(t, configMap["cohort_feature"])

	// cohort handler, match filter
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{}})
	configMap, _ = abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.Equal(t, float64(2), configMap["cohort_feature"]["val"].(float64))
	configMap, _ = abClient.GetFeatureConfigs("10301", "decisionId", attributes)
	assert.Equal(t, float64(2), configMap["cohort_feature"]["val"].(float64))

	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"12345"}})
	configMap, _ = abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.Equal(t, float64(1), configMap["cohort_feature"]["val"])
	configMap, _ = abClient.GetFeatureConfigs("10301", "decisionId", attributes)
	assert.Equal(t, float64(1), configMap["cohort_feature"]["val"])

	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123456"}})
	configMap, _ = abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.Equal(t, float64(2), configMap["cohort_feature"]["val"])
	configMap, _ = abClient.GetFeatureConfigs("10301", "decisionId", attributes)
	assert.Equal(t, float64(2), configMap["cohort_feature"]["val"])

	// cohort handler, not match filter
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"54321"}})
	configMap, _ = abClient.GetAllFeatureConfigs("decisionId", attributes)
	assert.Nil(t, configMap["cohort_feature"])
	configMap, _ = abClient.GetFeatureConfigs("10301", "decisionId", attributes)
	assert.Nil(t, configMap["cohort_feature"])
}
