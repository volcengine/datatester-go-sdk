/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package httper

import "net/http"

type HResponse struct {
	Body    []byte
	Headers http.Header
	Code    int
}
