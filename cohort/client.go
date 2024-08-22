/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cohort

type Client interface {
	UseCohort(decisionId string, cohortIds []string) []string
}
