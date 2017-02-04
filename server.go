package letsrest

import (
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"net/http"
)

func NewServer(r Requester, s RequestStore) *Server {
	return &Server{requester: r, store: s}
}

type Server struct {
	requester Requester
	store     RequestStore
}

func IrisHandler(requester Requester, store RequestStore) *iris.Framework {
	srv := NewServer(requester, store)
	api := iris.New()
	api.UseFunc(logger.New())

	api.Get("/", func(ctx *iris.Context) {
		ctx.JSON(http.StatusOK, "OK")
		return
	})

	v1 := api.Party("/api/v1")
	{
		v1.Get("/", func(ctx *iris.Context) {
			ctx.JSON(http.StatusOK, "OK")
			return
		})

		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6111/users/*anything
		//api.OnError(404, userNotFoundHandler)

		// http://localhost:6111/users/42
		// Method: "GET"
		v1.Put("/requests", srv.CreateRequest)
		v1.Get("/requests/:id", srv.GetRequest)
	}

	api.Build()
	return api
}

func (s *Server) CreateRequest(ctx *iris.Context) {
	clientRequest := &ClientRequest{}

	if err := ctx.ReadJSON(clientRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	requestID, err := s.store.Save(clientRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	v := struct {
		ReqID string `json:"req_id"`
	}{
		ReqID: requestID,
	}

	ctx.JSON(http.StatusCreated, v)
}

func (s *Server) GetRequest(ctx *iris.Context) {
	cReq, err := s.store.Get(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if cReq == nil {
		ctx.JSON(http.StatusNotFound, RequestNotFoundResponse(ctx.Param("id")))
		return
	}

	ctx.JSON(http.StatusOK, cReq)
}
