/**
 * Apache 2.0
 * Copyright 2022 Beijing Volcano Engine Technology Co., Ltd.
 */

package httper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/volcengine/datatester-go-sdk/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
	DefaultRetries = 3
)

type Requester interface {
	Get(url string) (*HResponse, error)
	Post(url string, body interface{}) (*HResponse, error)
}

type Header struct {
	Name, Value string
}

func WithTimeout(timeout time.Duration) func(r *HRequester) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return func(r *HRequester) {
		r.client = http.Client{
			Timeout: timeout,
			// skip ssl verify
			//Transport: &http.Transport{
			//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//},
		}
	}
}

func WithRetries(retries int) func(r *HRequester) {
	if retries <= 0 {
		retries = DefaultRetries
	}
	return func(r *HRequester) {
		r.retries = retries
	}
}

func WithHeaders(headers ...Header) func(r *HRequester) {
	return func(r *HRequester) {
		r.headers = []Header{}
		r.headers = append(r.headers, headers...)
	}
}

type HRequester struct {
	client  http.Client
	retries int
	headers []Header
}

func NewHTTPRequester(options ...func(*HRequester)) *HRequester {
	requester := HRequester{
		retries: DefaultRetries,
		headers: []Header{{
			"Content-Type", "application/json"}, {"Accept", "application/json"}},
		client: http.Client{
			Timeout: DefaultTimeout,
			// skip ssl verify
			//Transport: &http.Transport{
			//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//},
		},
	}

	for _, fn := range options {
		fn(&requester)
	}
	return &requester
}

func (r *HRequester) Get(url string) (*HResponse, error) {
	return r.Do(url, "GET", nil)
}

func (r *HRequester) Post(url string, body interface{}) (*HResponse, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return r.Do(url, "POST", bytes.NewBuffer(b))
}

func streamToString(stream io.Reader) string {
	if stream == nil {
		return ""
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (r *HRequester) Do(url, method string, body io.Reader) (response *HResponse, err error) {
	send := func(request *http.Request) (response *HResponse, e error) {
		resp, doErr := r.client.Do(request)
		if doErr != nil {
			log.WithFields(
				log.Fields{
					"requestHost": request.Host,
					"requestBody": streamToString(request.Body),
					"error":       doErr}).Error("failed to send request")
			return nil, doErr
		}
		defer func() {
			if e := resp.Body.Close(); e != nil {
				log.WithFields(log.Fields{
					"request.url": request.URL,
					"error":       e,
				}).Warn("can't close body")
			}
		}()

		if resp.StatusCode >= http.StatusBadRequest {
			log.WithFields(log.Fields{"status code": resp.StatusCode}).Error("error status code")
			return nil, errors.New(resp.Status)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("failed to read body")
			return nil, err
		}

		return &HResponse{Body: body, Headers: resp.Header, Code: resp.StatusCode}, nil
	}

	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		log.WithFields(log.Fields{"url": url, "error": reqErr}).Error("failed to make request")
		return nil, reqErr
	}

	// add headers
	for _, h := range r.headers {
		req.Header.Add(h.Name, h.Value)
	}

	for i := 0; i < r.retries; i++ {
		if response, err = send(req); err == nil {
			triedMsg := fmt.Sprintf("tried %d time(s)", i+1)
			log.WithFields(log.Fields{"url": url, "triedMsg": triedMsg}).Debug("retry fetch...")
			return response, nil
		}
		if i != r.retries {
			delay := time.Duration(500) * time.Millisecond
			time.Sleep(delay)
		}
	}
	return response, err
}

var _ Requester = &HRequester{}
