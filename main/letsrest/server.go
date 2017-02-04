package main

import (
	"errors"
	"github.com/itimofeev/letsrest"
	"github.com/kataras/iris"
	"net/http"
)

func NewServer(r Requester) *Server {
	return &Server{r:r}
}

type Server struct {
	r Requester
}

func (s *Server) CreateRequest(ctx *iris.Context) {
	clientRequest := &letsrest.ClientRequest{}

	if err := ctx.ReadJSON(clientRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := validateRequest(clientRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	clientResponse, err := s.MakeExternalRequest(clientRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, clientResponse)
}

func validateRequest(request *letsrest.ClientRequest) error {
	if request.URL == "" {
		return errors.New("Empty url")
	}

	return nil
}

func (s *Server) MakeExternalRequest(cReq *letsrest.ClientRequest) (*letsrest.ClientResponse, error) {
	resp, err := s.r.Do(cReq)
	if err != nil {
		return nil, err
	}

	cResp := &letsrest.ClientResponse{StatusCode: resp.StatusCode}

	return cResp, nil
}

