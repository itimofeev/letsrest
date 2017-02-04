package letsrest

import (
	"errors"
)

type Requester interface {
	Do(request *ClientRequest) (*ClientResponse, error)
}

type HTTPRequester struct {

}

func (r*HTTPRequester) Do(request *ClientRequest) (cResp *ClientResponse,err error){
		//req, err := http.NewRequest(cReq.Method, cReq.URL, nil)
	//if err != nil {
	//	return nil, err
	//}
	return nil, errors.New("Not implemented")
}

