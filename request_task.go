package letsrest

import (
	"math"
	"strings"
	"time"
)

// информация задаче на выполнение запроса
type Request struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	RequestData *RequestData `json:"data,omitempty" bson:"data"`
	Response    *Response    `json:"response,omitempty" bson:"response"`
	Status      *ExecStatus  `json:"status" bson:"status"`

	UserID string `json:"user_id" bson:"user_id"`
}

type RequestData struct {
	URL     string   `json:"url" bson:"url"`
	Method  string   `json:"method" bson:"method"`
	Headers []Header `json:"headers" bson:"headers"`
	Body    string   `json:"body" bson:"body"`
}

type HeaderSlice []Header

func (p HeaderSlice) Len() int           { return len(p) }
func (p HeaderSlice) Less(i, j int) bool { return strings.Compare(p[i].Name, p[j].Name) < 0 }
func (p HeaderSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// информация об ответе на запрос
type Response struct {
	StatusCode  int           `json:"status_code,omitempty" bson:"status_code"`
	Headers     HeaderSlice   `json:"headers,omitempty" bson:"headers"`
	Body        string        `json:"body" bson:"body"`
	BodyLen     int           `json:"body_len,omitempty" bson:"body_len"`
	ContentType string        `json:"content_type,omitempty" bson:"content_type"`
	Duration    time.Duration `json:"duration,omitempty" bson:"duration"` // in ns
}

// Header данные о заголовке
type Header struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type ExecStatus struct {
	Status string `json:"status" bson:"status"`
	Error  string `json:"error,omitempty" bson:"error"`
}

type User struct {
	ID           string `json:"id" bson:"_id"`
	RequestLimit int    `json:"request_limit" bson:"request_limit"`
}

func (u *User) GetRequestLimit() int {
	if u.RequestLimit == 0 {
		return math.MaxInt32
	}
	return u.RequestLimit
}
