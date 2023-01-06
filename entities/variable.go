/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package entities

type Variable struct {
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
