package letsrest

import (
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
	req, err := http.NewRequest(request.Method, request.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	var h []Header
	for key, value := range resp.Header {
		h = append(h, Header{Name: key, Value: strings.Join(value, ", ")})
	}

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cResp = &Response{
		ID:         request.ID,
		StatusCode: resp.StatusCode,
		Headers:    h,
		Body:       bodyData,
		BodyLen:    len(bodyData),
	}
	return
}
