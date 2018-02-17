package letsrest

import (
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Requester interface {
	Do(request *RequestData) (*Response, error)
}

const defaultBodyLimit int64 = 1024 * 1024 * 10 // 10MB

func NewHTTPRequester(maxBodySize ...int64) *HTTPRequester {
	limit := defaultBodyLimit
	if len(maxBodySize) > 0 {
		limit = maxBodySize[0]
	}
	return &HTTPRequester{
		maxBodySize: limit,
		client:      newHTTPClient(),
	}
}

func newHTTPClient() *HTTPClientDefault {
	return &HTTPClientDefault{
		client: http.DefaultClient,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClientDefault struct {
	client *http.Client
}

func (c *HTTPClientDefault) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

type HTTPRequester struct {
	maxBodySize int64
	client      HTTPClient
}

func (r *HTTPRequester) Do(request *RequestData) (cResp *Response, err error) {
	var reader io.Reader
	if len(request.Body) > 0 {
		reader = strings.NewReader(request.Body)
	}

	req, err := http.NewRequest(request.Method, request.URL, reader)
	if err != nil {
		return nil, err
	}

	for _, header := range request.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	start := time.Now()

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	var h HeaderSlice
	for key, value := range resp.Header {
		h = append(h, Header{Name: key, Value: strings.Join(value, ", ")})
	}
	sort.Sort(h)

	contentTypeHeader := findHeader("Content-Type", h)
	contentType := ""
	if contentTypeHeader != nil {
		contentType = contentTypeHeader.Value
	}

	limitedReader := &LimitedErrReader{N: r.maxBodySize, R: resp.Body}
	bodyData, err := ioutil.ReadAll(limitedReader)
	cResp = &Response{
		StatusCode:  resp.StatusCode,
		Headers:     h,
		BodyLen:     len(bodyData),
		Body:        string(bodyData),
		ContentType: contentType,
		Duration:    time.Now().Sub(start),
	}
	return
}
