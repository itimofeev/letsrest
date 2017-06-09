package letsrest

import (
	"errors"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var store = NewRequestStore(&testRequester{})

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
	request := createRequest(t)

	getResp := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		Expect().
		Status(http.StatusOK).
		JSON()

	getResp.Object().ValueEqual("id", request.ID)
	getResp.Object().Value("status").Object().ValueEqual("status", "idle")
}

func TestServer_GetNotExistedRequest(t *testing.T) {
	v := tester(t).GET("/api/v1/requests/{ID}", "someNotExistedID").
		Expect().
		Status(http.StatusNotFound).
		JSON()

	v.Object().Value("key").Equal(ReqNotFoundKey)
	v.Object().ValueEqual("params", Params{"id": "someNotExistedID"})
}

func TestServer_GetReadyResponse(t *testing.T) {
	request := createRequest(t)

	resp := &Response{StatusCode: 200}
	store.SetResponse(request.ID, resp, nil)

	obj := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("response").Equal(resp)
	obj.Value("status").Object().ValueEqual("status", "done")
}

func TestServer_GetErrorResponse(t *testing.T) {
	request := createRequest(t)

	store.SetResponse(request.ID, nil, errors.New("error.while.do.request"))

	obj := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("status").Object().ValueEqual("status", "error")
	obj.Value("status").Object().ValueEqual("error", "error.while.do.request")
}

func TestServer_GetNotReadyResponse(t *testing.T) {
	request := createRequest(t)

	r := tester(t).GET("/api/v1/requests/{ID}", request.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	r.Value("status").Object().ValueEqual("status", "idle")
}

func TestServer_ExecRequest(t *testing.T) {
	request := createRequest(t)

	data := RequestData{Method: "hello", URL: "there"}

	r := tester(t).PUT("/api/v1/requests/{ID}", request.ID).
		WithJSON(data).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	r.Value("status").Object().ValueEqual("status", "in_progress")
	r.ValueEqual("data", data)
}

func createRequest(t *testing.T) *Request {
	request := &Request{Name: "some name"}

	resp := tester(t).POST("/api/v1/requests").
		WithJSON(request).
		Expect().
		Status(http.StatusCreated).
		JSON()

	resp.Object().ContainsKey("id")
	request.ID = resp.Object().Value("id").Raw().(string)

	return request
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
