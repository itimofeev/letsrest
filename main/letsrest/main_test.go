package main

import (
	"github.com/itimofeev/letsrest"
	"net/http"
	"testing"
)

type testRequester struct {
}

func (r *testRequester) Do(request *letsrest.ClientRequest) (*letsrest.ClientResponse, error) {
	return nil, nil
}

func TestSimpleApiCalls(t *testing.T) {
	e := IrisHandler(&testRequester{}).Tester(t)

	//cReq := letsrest.ClientRequest{URL: "http://google.com", Method: "DELETE"}
	//
	//e.PUT("/api/v1/requests").
	//	WithJSON(cReq).
	//	Expect().
	//	Status(http.StatusOK)

	e.GET("/").
		Expect().
		Status(http.StatusOK)

	e.GET("/api/v1").
		Expect().
		Status(http.StatusOK)
}

//func irisTester(t *testing.T) *httpexpect.Expect {
//	handler := IrisHandler()
//
//	return httpexpect.WithConfig(httpexpect.Config{
//		BaseURL: "http://example.com",
//		Client: &http.Client{
//			Transport: httpexpect.NewBinder(handler.Router),
//			Jar:       httpexpect.NewJar(),
//		},
//		Reporter: httpexpect.NewAssertReporter(t),
//		Printers: []httpexpect.Printer{
//			httpexpect.NewDebugPrinter(t, true),
//		},
//	})
//}
