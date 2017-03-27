package letsrest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requester interface {
	Do(request *RequestTask) (*Response, error)
}

func NewHTTPRequester() *HTTPRequester {
	return &HTTPRequester{}
}

type HTTPRequester struct {
}

func (r *HTTPRequester) Do(request *RequestTask) (cResp *Response, err error) {
	var reader io.Reader
	if len(request.Body) > 0 {
		reader = bytes.NewReader(request.Body)
	}

	req, err := http.NewRequest(request.Method, request.URL, reader)
	if err != nil {
		return nil, err
	}

	for _, header := range request.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var h []Header
	for key, value := range resp.Header {
		h = append(h, Header{Name: key, Value: strings.Join(value, ", ")})
	}

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
		ID:          request.ID,
		StatusCode:  resp.StatusCode,
		Headers:     h,
		Body:        bodyData,
		BodyLen:     len(bodyData),
		ContentType: contentType,
	}
	return
}
