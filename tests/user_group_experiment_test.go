package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/volcengine/datatester-go-sdk/consts"
	"testing"
)

func TestUserGroupExperiment(t *testing.T) {
	abClient, _ := setUp()
	attributes := make(map[string]interface{})

	// not match user group rule
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{}})
	result, _ := abClient.Activate("user_group_config", "decisionId",
		"trackId", nil, attributes)
	assert.Nil(t, result)
	assert.Nil(t, attributes[consts.UserGroupRelation])

	// hit one user group rule
	attributes = make(map[string]interface{})
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123"}})
	result, _ = abClient.Activate("user_group_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	assert.Equal(t, "100", attributes[consts.UserGroupRelation].(map[string]string)["120340"])

	attributes = make(map[string]interface{})
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"456"}})
	result, _ = abClient.Activate("user_group_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	assert.Equal(t, "101", attributes[consts.UserGroupRelation].(map[string]string)["120340"])

	// hit multiple user group rules
	// mmh3.hash('decisionId:1fe7472a74721b00db5c97e1e6764b23:user-group') % 2 == 0
	attributes = make(map[string]interface{})
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123", "456"}})
	result, _ = abClient.Activate("user_group_config", "decisionId",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	assert.Equal(t, "100", attributes[consts.UserGroupRelation].(map[string]string)["120340"])

	// mmh3.hash('decisionId2:1fe7472a74721b00db5c97e1e6764b23:user-group') % 2 == 1
	attributes = make(map[string]interface{})
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{"123", "456"}})
	result, _ = abClient.Activate("user_group_config", "decisionId2",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	assert.Equal(t, "101", attributes[consts.UserGroupRelation].(map[string]string)["120340"])

	// test user not filter allow list
	// mmh3.hash('user-group-exp-test-user:1fe7472a74721b00db5c97e1e6764b23:user-group') % 2 == 0
	attributes = make(map[string]interface{})
	abClient.SetCohortHandler(&fakeCohortClient{fakeCohortIds: []string{}})
	result, _ = abClient.Activate("user_group_config", "user-group-exp-test-user",
		"trackId", nil, attributes)
	assert.True(t, result.(bool))
	assert.Equal(t, "100", attributes[consts.UserGroupRelation].(map[string]string)["120340"])
}
