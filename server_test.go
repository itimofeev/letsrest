package letsrest

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var requestID = "hello"
var storedReq *ClientRequest

type testRequester struct {
}

func (r *testRequester) Do(request *ClientRequest) (*ClientResponse, error) {
	return nil, nil
}

type testRequestStore struct {
}

func (s testRequestStore) Save(cReq *ClientRequest) (*ClientRequest, error) {
	cReq.ID = requestID
	return cReq, nil
}
func (s testRequestStore) Get(id string) (*ClientRequest, error) {
	return storedReq, nil
}
func (s testRequestStore) Delete(id string) error {
	return nil
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
	cReq := ClientRequest{URL: "http://somedomain.com", Method: "POST"}

	resp := tester(t).PUT("/api/v1/requests").
		WithJSON(cReq).
		Expect().
		Status(http.StatusCreated).
		JSON()

	reqIDValue := resp.Object().Value("id")
	reqIDValue.Equal(requestID)
	cReq.ID = reqIDValue.Raw().(string)

	storedReq = &cReq

	getResp := tester(t).GET("/api/v1/requests/{reqID}", cReq.ID).
		Expect().
		Status(http.StatusOK).
		JSON()

	getResp.Object().Equal(cReq)
	getResp.Object().ValueEqual("id", cReq.ID)
}

func TestServer_GetNotExistedRequest(t *testing.T) {
	storedReq = nil
	v := tester(t).GET("/api/v1/requests/{reqID}", "someNotExistedID").
		Expect().
		Status(http.StatusNotFound).
		JSON()

	v.Object().Value("key").Equal(ReqNotFoundKey)
	v.Object().ValueEqual("params", Params{"id": "someNotExistedID"})
}

func tester(t *testing.T) *httpexpect.Expect {
	return IrisHandler(&testRequester{}, testRequestStore{}).Tester(t)
}
