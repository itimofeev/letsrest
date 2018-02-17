package letsrest

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type fakeHTTPClient struct {
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader("more than 10 length")),
	}, nil
}

func Test_Requester_ShouldLimitBody(t *testing.T) {
	requester := HTTPRequester{
		maxBodySize: 10,
		client:      &fakeHTTPClient{},
	}

	reqData := &RequestData{}

	resp, err := requester.Do(reqData)

	assert.EqualError(t, err, bodySizeLimitExceededErr.Error())
	assert.Len(t, resp.Body, 10)
}
