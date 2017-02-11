package letsrest

import (
	"net/http"
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

	cResp = &Response{
		ID:         request.ID,
		StatusCode: resp.StatusCode,
	}
	return
}
