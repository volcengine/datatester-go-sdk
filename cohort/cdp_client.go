/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package cohort

import (
	"context"
	cdpCli "github.com/volcengine/cdp-openapi-sdk-go"
	"github.com/volcengine/datatester-go-sdk/log"
	"strings"
)

type CdpClient struct {
	cdpApi   *cdpCli.APIClient
	tenantId int64
}

func NewCdpClient(cdpApi *cdpCli.APIClient, tenantId int64) *CdpClient {
	return &CdpClient{cdpApi: cdpApi, tenantId: tenantId}
}

func (c *CdpClient) UseCohort(decisionId string, cohortIds []string) []string {
	if c.cdpApi == nil || len(cohortIds) == 0 {
		return nil
	}

	// gen req
	var cohortIdStrBuilder strings.Builder
	cohortIdStrBuilder.WriteString(cohortIds[0])
	for _, id := range cohortIds[1:] {
		cohortIdStrBuilder.WriteString(",")
		cohortIdStrBuilder.WriteString(id)
	}

	// call api
	resp, _, err := c.cdpApi.OnlineApi.QueryUserSeg(context.Background(), decisionId, cohortIdStrBuilder.String(), c.tenantId)
	if err != nil {
		log.ErrorF("call cdp cohort api failed, error = %v", err.Error())
		return nil
	}

	return resp.Data
}
