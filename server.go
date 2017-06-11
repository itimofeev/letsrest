package letsrest

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/iris"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var secretForJwt = []byte("123") // TODO secured string

var log = logrus.New()

func NewServer(s RequestStore) *Server {
	log.Out = os.Stdout
	formatter := new(logrus.TextFormatter)
	formatter.ForceColors = true
	log.Formatter = formatter
	log.Level = logrus.DebugLevel

	//anonymLimiter := rate.NewLimiter(rate.Every(time.Duration(200)*time.Millisecond), 5)

	return &Server{store: s}
}

type Server struct {
	store         RequestStore
	anonymLimiter *rate.Limiter
}

func IrisHandler(store RequestStore) *iris.Framework {
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

	v1 := api.Party("/api/v1", srv.CheckAuthToken)
	{
		v1.Get("/", func(ctx *iris.Context) {
			ctx.JSON(http.StatusOK, "OK")
			return
		})

		// Fire userNotFoundHandler when Not Found
		// inside http://localhost:6111/users/*anything
		//api.OnError(404, userNotFoundHandler)

		v1.Get("/authTokens", srv.CreateAuthToken)

		requests := v1.Party("/requests")

		requests.Post("", srv.CreateRequest)
		requests.Put("/:id", srv.ExecRequest)
		requests.Get("/:id", srv.GetRequest)
		requests.Get("", srv.GetRequests)

		v1.Any("/test", srv.Test)
	}

	return api
}

func (s *Server) CheckAuthToken(ctx *iris.Context) {
	_, ok := ctx.Request.Header["Authorization"]

	if !ok {
		ctx.ResponseWriter.Header().Add("Content-Type", "application/json")
		ctx.ResponseWriter.WriteHeader(http.StatusTooManyRequests)
		ctx.ResponseWriter.Write([]byte("You have reached maximum request limit."))
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// LetsRestClaims claims in terminology of jwt just a data that serialized in jwt token
type LetsRestClaims struct {
	UserID string
	jwt.StandardClaims
}

type Auth struct {
	AuthToken string `json:"auth_token"`
}

func (s *Server) CreateAuthToken(ctx *iris.Context) {
	userID, err := uuid.NewV4()
	Must(err, "uuid.NewV4()")
	expDate := time.Now().Add(time.Hour * 24 * 355 * 10) // 10 years

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, LetsRestClaims{
		UserID: userID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expDate.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LetsRest",
		},
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretForJwt))
	Must(err, "token.SignedString([]byte(secretForJwt))")

	ctx.JSON(http.StatusOK, Auth{AuthToken: tokenString})
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
	request, err := s.store.CreateRequest(name.Name)
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
	requests, err := s.store.List()
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
