/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package handler

type UserAbInfoHandler interface {
	Query(decisionID string) string
	CreateOrUpdate(decisionID, experiment2Variant string) bool
	NeedPersistData() bool
}

type DefaultUserAbInfoHandler struct{}

func (u *DefaultUserAbInfoHandler) Query(decisionID string) string {
	return ""
}

func (u *DefaultUserAbInfoHandler) CreateOrUpdate(decisionID, experiment2Variant string) bool {
	return true
}

func (u *DefaultUserAbInfoHandler) NeedPersistData() bool {
	return false
}

func NewDefaultUserAbInfoHandler() *DefaultUserAbInfoHandler {
	return &DefaultUserAbInfoHandler{}
}
