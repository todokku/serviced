// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 7.7.0: DO NOT EDIT

package esapi

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func newMLEstimateModelMemoryFunc(t Transport) MLEstimateModelMemory {
	return func(body io.Reader, o ...func(*MLEstimateModelMemoryRequest)) (*Response, error) {
		var r = MLEstimateModelMemoryRequest{Body: body}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// MLEstimateModelMemory - Estimates the model memory
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-apis.html.
//
type MLEstimateModelMemory func(body io.Reader, o ...func(*MLEstimateModelMemoryRequest)) (*Response, error)

// MLEstimateModelMemoryRequest configures the ML Estimate Model Memory API request.
//
type MLEstimateModelMemoryRequest struct {
	Body io.Reader

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r MLEstimateModelMemoryRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "POST"

	path.Grow(len("/_ml/anomaly_detectors/_estimate_model_memory"))
	path.WriteString("/_ml/anomaly_detectors/_estimate_model_memory")

	params = make(map[string]string)

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), r.Body)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if r.Body != nil {
		req.Header[headerContentType] = headerContentTypeJSON
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f MLEstimateModelMemory) WithContext(v context.Context) func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f MLEstimateModelMemory) WithPretty() func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f MLEstimateModelMemory) WithHuman() func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f MLEstimateModelMemory) WithErrorTrace() func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f MLEstimateModelMemory) WithFilterPath(v ...string) func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f MLEstimateModelMemory) WithHeader(h map[string]string) func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f MLEstimateModelMemory) WithOpaqueID(s string) func(*MLEstimateModelMemoryRequest) {
	return func(r *MLEstimateModelMemoryRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
