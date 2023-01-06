/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package event

import (
	"github.com/volcengine/datatester-go-sdk/event/model"
	"github.com/volcengine/datatester-go-sdk/event/util"
	"github.com/volcengine/datatester-go-sdk/log"
	"strconv"
)

type Dispatcher interface {
	DispatchEvent(trackId, vid string, attributes map[string]interface{}) error
	DispatchEventWithIdType(trackId, vid, uuidType string, attributes map[string]interface{}) error
}

type DispatcherImpl struct {
	appId  uint32
	header *model.Header
}

func NewDispatcher(token string) *DispatcherImpl {
	appId, err := strconv.ParseUint(token, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("token is invalid")
		return nil
	}
	d := &DispatcherImpl{
		appId:  uint32(appId),
		header: util.CreateHeader(uint32(appId), util.GetTimezone()),
	}
	return d
}

func (d *DispatcherImpl) DispatchEvent(trackId, vid string, attributes map[string]interface{}) error {
	return d.DispatchEventWithIdType(trackId, vid, "", attributes)
}

func (d *DispatcherImpl) DispatchEventWithIdType(trackId, vid, uuidType string,
	attributes map[string]interface{}) error {
	e := util.CreateExposureEvent(vid, uuidType)
	user := util.CreateUser(trackId, uuidType, attributes, GetAnonymousUserConfig().Enable,
		GetAnonymousUserConfig().IsSaas)
	events := &model.ExposureEvents{
		User:   user,
		Header: d.header,
		Events: []*model.Event{e},
	}
	if err := GetInstance().CollectEvents(events); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("events dispatcher failed")
		return err
	}
	return nil
}

type DispatcherImpl4Test struct {
}

func NewDispatcher4Test() *DispatcherImpl4Test {
	d := &DispatcherImpl4Test{}
	return d
}

func (d *DispatcherImpl4Test) DispatchEvent(trackId, vid string, attributes map[string]interface{}) error {
	return nil
}

func (d *DispatcherImpl4Test) DispatchEventWithIdType(trackId, vid, uuidType string,
	attributes map[string]interface{}) error {
	return nil
}
