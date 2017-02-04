package letsrest

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var requestID = "hello"
var cReq *ClientRequest

type testRequester struct {
}

func (r *testRequester) Do(request *ClientRequest) (*ClientResponse, error) {
	return nil, nil
}

type testRequestStore struct {
}

func (s testRequestStore) Save(*ClientRequest) (string, error) {
	return requestID, nil
}
func (s testRequestStore) Get(id string) (*ClientRequest, error) {
	return nil, nil
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

	reqID := resp.Object().Value("req_id")
	reqID.Equal(requestID)

	//getResp := tester(t).GET("/api/v1/requests/{reqID}", reqID.String()).
	//	WithJSON(cReq).
	//	Expect().
	//	Status(http.StatusOK).
	//	JSON()
}

func TestServer_GetNotExistedRequest(t *testing.T) {
	v := tester(t).GET("/api/v1/requests/{reqID}", "someNotExistedID").
		WithJSON(cReq).
		Expect().
		Status(http.StatusNotFound).
		JSON()

	v.Object().Value("key").Equal(ReqNotFoundKey)
	v.Object().ValueEqual("params", Params{"id": "someNotExistedID"})
}

func tester(t *testing.T) *httpexpect.Expect {
	return IrisHandler(&testRequester{}, testRequestStore{}).Tester(t)
}
