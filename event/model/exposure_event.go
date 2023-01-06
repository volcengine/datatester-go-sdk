/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package model

import "fmt"

type ExposureEvents struct {
	User   *User    `json:"user,omitempty"`
	Header *Header  `json:"header,omitempty"`
	Events []*Event `json:"events,omitempty"`
}

func (v *ExposureEvents) String() string {
	return fmt.Sprintf("{user: %v, header: %v, events: %v}", v.User, v.Header, v.Events)
}

type User struct {
	UserUniqueId *string `json:"user_unique_id,omitempty"`
	UuidType     *string `json:"user_unique_id_type,omitempty"`
	DeviceId     *int64  `json:"device_id,omitempty"`
	BdDid        *string `json:"bddid,omitempty"`
	WebId        *int64  `json:"web_id,omitempty"`
}

func (v *User) String() string {
	return fmt.Sprintf("{user_unique_id: %v, user_unique_id_type: %v}", "***", *v.UuidType)
}

type Header struct {
	AppId    *uint32 `json:"app_id,omitempty"`
	Timezone *int32  `json:"timezone,omitempty"`
}

func (v *Header) String() string {
	return fmt.Sprintf("{app_id: %v, timezone: %v}", *v.AppId, *v.Timezone)
}

type Event struct {
	Event        *string `json:"event,omitempty"`
	Params       *string `json:"params,omitempty"`
	LocalTimeMs  *uint64 `json:"local_time_ms,omitempty"`
	AbSdkVersion *string `json:"ab_sdk_version,omitempty"`
}

func (v *Event) String() string {
	return fmt.Sprintf("{event: %v, params: %v, local_time_ms: %v, ab_sdk_version: %v}",
		*v.Event, *v.Params, *v.LocalTimeMs, *v.AbSdkVersion)
}
