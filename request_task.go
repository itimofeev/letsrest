package letsrest

// информация задаче на выполнение запроса
type Bucket struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Request  *Request    `json:"request,omitempty"`
	Response *Response   `json:"response,omitempty"`
	Status   *ExecStatus `json:"status"`
}

type Request struct {
	URL     string   `json:"url"`
	Method  string   `json:"method"`
	Headers []Header `json:"headers"`
}

// информация об ответе на запрос
type Response struct {
	StatusCode  int      `json:"status_code,omitempty"`
	Headers     []Header `json:"headers,omitempty"`
	BodyLen     int      `json:"body_len,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
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
