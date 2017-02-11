package letsrest

import (
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"net/http"
)

func NewServer(r Requester, s RequestStore) *Server {
	return &Server{requester: r, store: s, taskCh: make(chan *RequestTask, 100)}
}

type Server struct {
	requester Requester
	store     RequestStore
	taskCh    chan *RequestTask
}

func IrisHandler(requester Requester, store RequestStore) (*iris.Framework, *Server) {
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

		v1.Put("/requests", srv.CreateRequest)
		v1.Get("/requests/:id", srv.GetRequest)
		v1.Get("/requests/:id/responses", srv.GetResponse)

		v1.Get("/test", srv.Test)
	}

	return api, srv
}

func (s *Server) ListenForTasks() {
	for task := range s.taskCh {
		resp, err := s.requester.Do(task)
		s.store.SetResponse(task.ID, resp, err)
	}
}

func (s *Server) CreateRequest(ctx *iris.Context) {
	requestTask := &RequestTask{}

	if err := ctx.ReadJSON(requestTask); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	requestTask, err := s.store.Save(requestTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	s.taskCh <- requestTask

	ctx.JSON(http.StatusCreated, requestTask)
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

func (s *Server) GetResponse(ctx *iris.Context) {
	cResp, err := s.store.GetResponse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if cResp == nil {
		ctx.JSON(http.StatusNotFound, RequestNotFoundResponse(ctx.Param("id")))
		return
	}

	if cResp.Status.Status == "in_progress" {
		ctx.JSON(http.StatusPartialContent, cResp)
		return
	}

	ctx.JSON(http.StatusOK, cResp)
}

func (s *Server) Test(ctx *iris.Context) {
	ctx.JSON(http.StatusOK, ctx.Request.URL.String())
}
