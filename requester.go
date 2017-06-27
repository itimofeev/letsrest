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

func NewHTTPRequester() *HTTPRequester {
	return &HTTPRequester{}
}

type HTTPRequester struct {
}

func (r *HTTPRequester) Do(request *RequestData) (cResp *Response, err error) {
	var reader io.Reader
	//if len(request.Body) > 0 {
	//	reader = bytes.NewReader(request.Body)
	//}

	req, err := http.NewRequest(request.Method, request.URL, reader)
	if err != nil {
		return nil, err
	}

	for _, header := range request.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	start := time.Now()

	resp, err := http.DefaultClient.Do(req)

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

	// TODO ограничение на размер ответа
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cResp = &Response{
		StatusCode:  resp.StatusCode,
		Headers:     h,
		BodyLen:     len(bodyData),
		ContentType: contentType,
		Duration:    time.Now().Sub(start),
	}
	return
}
