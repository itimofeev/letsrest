package main

import (
	"errors"
	"github.com/itimofeev/letsrest"
)

type Requester interface {
	Do(request *letsrest.ClientRequest) (*letsrest.ClientResponse, error)
}

type HTTPRequester struct {

}

func (r*HTTPRequester) Do(request *letsrest.ClientRequest) (cResp *letsrest.ClientResponse,err error){
		//req, err := http.NewRequest(cReq.Method, cReq.URL, nil)
	//if err != nil {
	//	return nil, err
	//}
	return nil, errors.New("Not implemented")
}

