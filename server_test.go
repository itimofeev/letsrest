package letsrest

import (
	"encoding/json"
	"errors"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
	"net/http"
	"testing"
	"time"
)

var store = NewDataStore(&Config{MongoURL: "192.168.99.100:27017"}, NewWorkerPool(&testRequester{}))

type testRequester struct {
}

func (r *testRequester) Do(request *RequestData) (*Response, error) {
	return nil, nil
}

func TestServer_SimpleApiCalls(t *testing.T) {
	tester(t).GET("/").
		Expect().
		Status(http.StatusOK)

	tester(t).GET("/api/v1").
		Expect().
		Status(http.StatusOK)
}

func TestServer_CreateRequest(t *testing.T) {
	request, auth := createRequest(t)

	getResp := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusOK).
		JSON()

	getResp.Object().ValueEqual("id", request.ID)
	getResp.Object().Value("status").Object().ValueEqual("status", "idle")
}

func TestServer_GetNotExistedRequest(t *testing.T) {
	_, auth := createRequest(t)

	v := tester(t).GET("/api/v1/requests/{ID}", "someNotExistedID").
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusNotFound).
		JSON()

	v.Object().Value("key").Equal(ReqNotFoundKey)
	v.Object().ValueEqual("params", Params{"id": "someNotExistedID"})
}

func TestServer_GetReadyResponse(t *testing.T) {
	request, auth := createRequest(t)

	resp := &Response{StatusCode: 200, Body: "someBody"}
	store.SetResponse(request.ID, resp, nil)

	obj := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("response").Equal(resp)
	obj.Value("status").Object().ValueEqual("status", "done")
	obj.Value("status").Object().ValueEqual("status", "done")
}

func TestServer_GetErrorResponse(t *testing.T) {
	request, auth := createRequest(t)

	store.SetResponse(request.ID, nil, errors.New("error.while.do.request"))

	obj := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("status").Object().ValueEqual("status", "error")
	obj.Value("status").Object().ValueEqual("error", "error.while.do.request")
}

func TestServer_GetNotReadyResponse(t *testing.T) {
	request, auth := createRequest(t)

	r := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	r.Value("status").Object().ValueEqual("status", "idle")
}

func TestServer_ExecRequest(t *testing.T) {
	request, auth := createRequest(t)

	data := RequestData{Method: "hello", URL: "there", Body: "someBody"}

	r := tester(t).PUT("/api/v1/requests/{ID}", request.ID).
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		WithJSON(data).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	r.Value("status").Object().ValueEqual("status", "in_progress")
	r.Value("data").Object().ValueEqual("method", data.Method)
	r.Value("data").Object().ValueEqual("url", data.URL)
	r.Value("data").Object().ValueEqual("body", data.Body)
	r.ValueEqual("data", data)
}

func TestServer_ExecRateLimit(t *testing.T) {
	auth := createAuth(t)
	expect := tester(t)

	for n := 0; n < 10; n++ {
		request := createRequestWithAuth(t, auth)
		data := RequestData{Method: "hello", URL: "there", Body: "someBody"}

		status := expect.PUT("/api/v1/requests/{ID}", request.ID).
			WithHeader("Authorization", "Bearer "+auth.AuthToken).
			WithJSON(data).
			Expect().Raw().StatusCode

		if status == http.StatusTooManyRequests {
			return
		}
	}
	assert.Fail(t, "rate limiter not triggered")
}

func Test_limiter(t *testing.T) {
	execLimiter := rate.NewLimiter(rate.Every(time.Second), 2)

	assert.True(t, execLimiter.Allow())
	assert.True(t, execLimiter.Allow())
	assert.False(t, execLimiter.Allow())
	assert.False(t, execLimiter.Allow())

	time.Sleep(time.Second)
	assert.True(t, execLimiter.Allow())
}

func TestServer_ListUserRequests(t *testing.T) {
	_, auth := createRequest(t)
	createRequest(t)

	tester(t).GET("/api/v1/requests").
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		Expect().
		Status(http.StatusOK).
		JSON().Array().Length().Equal(1)
}

func createAuth(t *testing.T) *Auth {
	userAndAuthToken := tester(t).POST("/api/v1/authTokens").
		Expect().
		Status(http.StatusOK).Body()
	auth := &Auth{}
	require.Nil(t, json.Unmarshal([]byte(userAndAuthToken.Raw()), auth))

	return auth
}

func createRequestWithAuth(t *testing.T, auth *Auth) *Request {
	request := &Request{Name: "some name"}

	resp := tester(t).POST("/api/v1/requests").
		WithHeader("Authorization", "Bearer "+auth.AuthToken).
		WithJSON(request).
		Expect().
		Status(http.StatusOK).
		JSON()

	resp.Object().ValueEqual("name", "some name")
	resp.Object().ContainsKey("id")
	request.ID = resp.Object().Value("id").Raw().(string)
	return request
}

func createRequest(t *testing.T) (*Request, *Auth) {
	auth := createAuth(t)
	request := createRequestWithAuth(t, auth)

	return request, auth
}

func tester(t *testing.T) *httpexpect.Expect {
	handler := IrisHandler(store)
	handler.Build()
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://example.com",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler.Router),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
