package letsrest

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
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

func IrisHandler(store DataStore) *iris.Application {
	srv := NewServer(store)
	api := iris.New()
	api.Use(func(ctx context.Context) {
		d, _ := httputil.DumpRequest(ctx.Request(), true)
		fmt.Println(string(d))
		ctx.Next()
	})

	api.Get("/", func(ctx context.Context) {
		ctx.JSON("OK")
	})

	v1 := api.Party("/api/v1")
	{
		v1.Get("/", func(ctx context.Context) {
			ctx.JSON("OK")
			return
		})

		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6111/users/*anything
		//api.OnError(404, userNotFoundHandler)

		v1.Post("/authTokens", srv.CreateAuthToken)

		requests := v1.Party("/requests", srv.CheckAuthToken)

		requests.Post("", srv.CreateRequest)
		requests.Put("/{id:string}", srv.ExecRequest)
		requests.Patch("/{id:string}", srv.EditRequest)
		requests.Delete("/{id:string}", srv.DeleteRequest)
		requests.Get("", srv.ListRequests)
		requests.Post("/{id:string}/copy", srv.CopyRequest)

		v1.Get("/requests/{id:string}", srv.GetRequest)

		v1.Any("/test", srv.Test)
	}

	return api
}

func (s *Server) CheckAuthToken(ctx context.Context) {
	authHeader, ok := ctx.Request().Header["Authorization"]

	if !ok {
		ctx.ResponseWriter().Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte("Auth header not found"))
		ctx.StopExecution()
		return
	}

	authToken := strings.Replace(authHeader[0], "Bearer ", "", 1)

	user, err := userFromAuthToken(authToken)
	if err != nil {
		ctx.ResponseWriter().Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte("Unable to decode auth token"))
		ctx.StopExecution()
		return
	}
	user, err = s.store.GetUser(user.ID)
	if err != nil {
		ctx.ResponseWriter().Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte("Error retrieving user from store"))
		ctx.StopExecution()
		return
	}

	if user == nil {
		ctx.ResponseWriter().Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter().WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter().Write([]byte("User not found in store"))
		ctx.StopExecution()
		return
	}

	ctx.Values().Set("LetsRestUser", user)
	ctx.Next()
}

func (s *Server) CreateAuthToken(ctx context.Context) {
	user := createUser()
	err := s.store.PutUser(user)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("PutUser returned error %s", err.Error()))
		return
	}
	auth := createAuthToken(user)
	ctx.JSON(auth)
}

func (s *Server) CreateRequest(ctx context.Context) {
	name := &struct {
		Name string `json:"name"`
	}{}
	err := ctx.ReadJSON(name)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(err.Error())
		return
	}

	request, err := s.store.CreateRequest(ctx.Values().Get("LetsRestUser").(*User), name.Name)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}
	ctx.JSON(request)
}

func (s *Server) ExecRequest(ctx context.Context) {
	data := &RequestData{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(err.Error())
		return
	}
	req, err := s.store.ExecRequest(ctx.Params().Get("id"), data)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}
	ctx.JSON(req)
}

func (s *Server) EditRequest(ctx context.Context) {
	name := &struct {
		Name string `json:"name"`
	}{}
	err := ctx.ReadJSON(name)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(err.Error())
		return
	}

	request, err := s.store.EditRequest(ctx.Params().Get("id"), name.Name)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}
	ctx.JSON(request)
}

func (s *Server) DeleteRequest(ctx context.Context) {
	err := s.store.Delete(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}
	ctx.JSON("OK")
}

func (s *Server) GetRequest(ctx context.Context) {
	req, err := s.store.GetRequest(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}

	if req == nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(RequestNotFoundResponse(ctx.Params().Get("id")))
		return
	}

	ctx.JSON(req)
}

func (s *Server) CopyRequest(ctx context.Context) {
	req, err := s.store.CopyRequest(ctx.Values().Get("LetsRestUser").(*User), ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err.Error())
		return
	}

	ctx.JSON(req)
}

func (s *Server) ListRequests(ctx context.Context) {
	requests, err := s.store.List(ctx.Values().Get("LetsRestUser").(*User))
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(err.Error())
		return
	}
	ctx.JSON(requests)
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

func (s *Server) Test(ctx context.Context) {
	dump, _ := httputil.DumpRequest(ctx.Request(), true)
	ctx.WriteString(string(dump))
}
