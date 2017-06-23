package letsrest

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var secretForJwt = []byte("12345678901234567890") // TODO move to settings

var log = logrus.New()

func NewServer(s DataStore) *Server {
	log.Out = os.Stdout
	formatter := new(logrus.TextFormatter)
	formatter.ForceColors = true
	log.Formatter = formatter
	log.Level = logrus.DebugLevel

	//anonymLimiter := rate.NewLimiter(rate.Every(time.Duration(200)*time.Millisecond), 5)

	return &Server{store: s}
}

type Server struct {
	store         DataStore
	anonymLimiter *rate.Limiter
}

func IrisHandler(store DataStore) *iris.Framework {
	srv := NewServer(store)
	api := iris.New()
	api.UseFunc(logger.New())
	api.UseFunc(func(ctx *iris.Context) {
		d, _ := httputil.DumpRequest(ctx.Request, true)
		fmt.Println(string(d))
		ctx.Next()
	})

	api.Get("/", func(ctx *iris.Context) {
		ctx.JSON(http.StatusOK, "OK")
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

		v1.Post("/authTokens", srv.CreateAuthToken)

		requests := v1.Party("/requests", srv.CheckAuthToken)

		requests.Post("", srv.CreateRequest)
		requests.Put("/:id", srv.ExecRequest)
		requests.Get("/:id", srv.GetRequest)
		requests.Get("", srv.GetRequests)

		v1.Any("/test", srv.Test)
	}

	return api
}

func (s *Server) CheckAuthToken(ctx *iris.Context) {
	authHeader, ok := ctx.Request.Header["Authorization"]

	if !ok {
		ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("Auth header not found"))
		ctx.StopExecution()
		return
	}

	authToken := strings.Replace(authHeader[0], "Bearer ", "", 1)

	user, err := userFromAuthToken(authToken)
	if err != nil {
		ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("Unable to decode auth token"))
		ctx.StopExecution()
		return
	}
	user, err = s.store.GetUser(user.ID)
	if err != nil {
		ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("Error retrieving user from store"))
		ctx.StopExecution()
		return
	}

	if user == nil {
		ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte("User not found in store"))
		ctx.StopExecution()
		return
	}

	ctx.Set("LetsRestUser", user)
	ctx.Next()
}

func (s *Server) CreateAuthToken(ctx *iris.Context) {
	user := createUser()
	err := s.store.PutUser(user)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("PutUser returned error %s", err.Error()))
		return
	}
	auth := createAuthToken(user)
	ctx.JSON(http.StatusOK, auth)
}

func (s *Server) CreateRequest(ctx *iris.Context) {
	name := &struct {
		Name string `json:"name"`
	}{}
	err := ctx.ReadJSON(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	request, err := s.store.CreateRequest(ctx.Get("LetsRestUser").(*User), name.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, request)
}

func (s *Server) ExecRequest(ctx *iris.Context) {
	data := &RequestData{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req, err := s.store.ExecRequest(ctx.Param("id"), data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, req)
}

func (s *Server) GetRequest(ctx *iris.Context) {
	req, err := s.store.Get(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if req == nil {
		ctx.JSON(http.StatusNotFound, RequestNotFoundResponse(ctx.Param("id")))
		return
	}

	ctx.JSON(http.StatusOK, req)
}

func (s *Server) GetRequests(ctx *iris.Context) {
	requests, err := s.store.List(ctx.Get("LetsRestUser").(*User))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, requests)
}

func findHeader(name string, headers []Header) *Header {
	loweredName := strings.ToLower(name)
	for _, header := range headers {
		if strings.ToLower(header.Name) == loweredName {
			return &header
		}
	}
	return nil
}

func (s *Server) Test(ctx *iris.Context) {
	dump, _ := httputil.DumpRequest(ctx.Request, true)
	ctx.WriteString(string(dump))
}
