package main

import (
	"github.com/itimofeev/letsrest"
	"github.com/kataras/iris"
	"net/http"
)

func main() {
	s := &Server{}

	api := iris.Party("/api/v1")
	{
		api.Get("/", func(ctx *iris.Context) {
			ctx.JSON(http.StatusOK, "OK")
		})

		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6111/users/*anything
		//api.OnError(404, userNotFoundHandler)

		// http://localhost:6111/users/42
		// Method: "GET"
		api.Put("/request", s.CreateRequest)
	}

	iris.Listen(":6111")
}

type Server struct {
}

func (s *Server) CreateRequest(ctx *iris.Context) {
	clientRequest := &letsrest.ClientRequest{}

	if err := ctx.ReadJSON(clientRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	clientResponse := letsrest.ClientResponse{StatusCode: http.StatusOK}
	ctx.JSON(http.StatusOK, clientResponse)
}
