package letsrest

import (
	"errors"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var store = NewRequestStore()

type testRequester struct {
}

func (r *testRequester) Do(request *Request) (*Response, error) {
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

func TestServer_CreateBucket(t *testing.T) {
	bucket := createBucket(t)

	getResp := tester(t).GET("/api/v1/buckets/{ID}", bucket.ID).
		Expect().
		Status(http.StatusOK).
		JSON()

	getResp.Object().ValueEqual("id", bucket.ID)
	getResp.Object().Value("status").Object().ValueEqual("status", "idle")
}

func TestServer_GetNotExistedRequest(t *testing.T) {
	v := tester(t).GET("/api/v1/buckets/{ID}", "someNotExistedID").
		Expect().
		Status(http.StatusNotFound).
		JSON()

	v.Object().Value("key").Equal(ReqNotFoundKey)
	v.Object().ValueEqual("params", Params{"id": "someNotExistedID"})
}

func TestServer_GetReadyResponse(t *testing.T) {
	bucket := createBucket(t)

	resp := &Response{StatusCode: 200}
	store.SetResponse(bucket.ID, resp, nil)

	obj := tester(t).GET("/api/v1/buckets/{ID}/responses", bucket.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("response").Equal(resp)
	obj.Value("status").Object().ValueEqual("status", "done")
}

func TestServer_GetErrorResponse(t *testing.T) {
	bucket := createBucket(t)

	store.SetResponse(bucket.ID, nil, errors.New("error.while.do.request"))

	obj := tester(t).GET("/api/v1/buckets/{ID}/responses", bucket.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	obj.Value("status").Object().ValueEqual("status", "error")
	obj.Value("status").Object().ValueEqual("error", "error.while.do.request")
}

func TestServer_GetNotReadyResponse(t *testing.T) {
	bucket := createBucket(t)

	r := tester(t).GET("/api/v1/buckets/{ID}/responses", bucket.ID).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	r.Value("status").Object().ValueEqual("status", "idle")
}

func createBucket(t *testing.T) *Bucket {
	bucket := &Bucket{Name: "some name"}

	resp := tester(t).POST("/api/v1/buckets").
		WithJSON(bucket).
		Expect().
		Status(http.StatusCreated).
		JSON()

	resp.Object().ContainsKey("id")
	bucket.ID = resp.Object().Value("id").Raw().(string)

	return bucket
}

func tester(t *testing.T) *httpexpect.Expect {
	handler, _ := IrisHandler(&testRequester{}, store)
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
