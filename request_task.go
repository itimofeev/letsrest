package letsrest

import (
	"strings"
	"time"
)

// информация задаче на выполнение запроса
type Request struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	RequestData *RequestData `json:"data,omitempty"`
	Response    *Response    `json:"response,omitempty"`
	Status      *ExecStatus  `json:"status"`

	UserID string `json:"user_id"`
}

type RequestData struct {
	URL     string   `json:"url"`
	Method  string   `json:"method"`
	Headers []Header `json:"headers"`
}

type HeaderSlice []Header

func (p HeaderSlice) Len() int           { return len(p) }
func (p HeaderSlice) Less(i, j int) bool { return strings.Compare(p[i].Name, p[j].Name) < 0 }
func (p HeaderSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// информация об ответе на запрос
type Response struct {
	StatusCode  int           `json:"status_code,omitempty"`
	Headers     HeaderSlice   `json:"headers,omitempty"`
	BodyLen     int           `json:"body_len,omitempty"`
	ContentType string        `json:"content_type,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"` // in ns
}

// Header данные о заголовке
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExecStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type User struct {
	ID string `json:"id"`
}
