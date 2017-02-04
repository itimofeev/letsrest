package main

import (
	"github.com/kataras/iris"
	"net/http"
)

func main() {
	framework := IrisHandler(&HTTPRequester{})
	framework.Listen(":6111")
}

func IrisHandler(r Requester) *iris.Framework {
	s := NewServer(r)
	api := iris.New()

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
		v1.Put("/requests", s.CreateRequest)
	}

	api.Build()
	return api
}
