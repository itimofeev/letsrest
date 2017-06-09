package letsrest

import (
	"github.com/Sirupsen/logrus"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"github.com/speps/go-hashids"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strings"
	"time"
)

var log = logrus.New()

func NewServer(r Requester, s RequestStore) *Server {
	log.Out = os.Stdout
	formatter := new(logrus.TextFormatter)
	formatter.ForceColors = true
	log.Formatter = formatter
	log.Level = logrus.DebugLevel

	//anonymLimiter := rate.NewLimiter(rate.Every(time.Duration(200)*time.Millisecond), 5)

	return &Server{requester: r, store: s, bucketCh: make(chan *Bucket, 100)}
}

type Server struct {
	requester     Requester
	store         RequestStore
	bucketCh      chan *Bucket
	anonymLimiter *rate.Limiter
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

		v1.Put("/authTokens", srv.CheckAuthToken)

		requests := v1.Party("/buckets")

		requests.Post("", srv.CreateBucket)
		requests.Post("/:id/requests", srv.CreateRequest)
		requests.Get("/:id", srv.GetRequest)
		requests.Get("/:id/responses", srv.GetResponse)
		requests.Get("", srv.GetRequestTaskList)

		v1.Get("/test", srv.Test)
	}

	return api, srv
}

func (s *Server) ListenForTasks() {
	for bucket := range s.bucketCh {
		resp, err := s.requester.Do(bucket.Request)
		s.store.SetResponse(bucket.ID, resp, err)
	}
}

func (s *Server) CheckAuthToken(ctx *iris.Context) {
	_, ok := ctx.Request.Header["Authorization"]
	// для неавторизированных пользователей проверяем, что лимит запросов не превышен
	if !ok {
		if !s.anonymLimiter.Allow() {
			ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
			ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
			ctx.ResponseWriter.Write([]byte("You have reached maximum request limit."))
			ctx.StopExecution()
			return
		}

		authToken, err := hashids.New().Encode([]int{time.Now().Second()})
		Must(err, "hashids.New().Encode([]int{time.Now().Second()})")
		ctx.ResponseWriter.Header().Add("X-LetsRest-AuthToken", authToken)
	}
	ctx.Next()
}

func (s *Server) CreateBucket(ctx *iris.Context) {
	name := &struct {
		Name string `json:"name"`
	}{}
	err := ctx.ReadJSON(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bucket, err := s.store.CreateBucket(name.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, bucket)
}


func (s *Server) CreateRequest(ctx *iris.Context) {
	req := &Request{}
	err := ctx.ReadJSON(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = s.store.CreateRequest(ctx.Param("id"), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, "")
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
	bucket, err := s.store.Get(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if bucket == nil {
		ctx.JSON(http.StatusNotFound, RequestNotFoundResponse(ctx.Param("id")))
		return
	}

	ctx.JSON(http.StatusOK, ResponseData{Response: bucket.Response, Status: bucket.Status})
}

func (s *Server) GetRequestTaskList(ctx *iris.Context) {
	taskList, err := s.store.List()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, taskList)
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
	ctx.JSON(http.StatusOK, ctx.Request.URL.String())
}
