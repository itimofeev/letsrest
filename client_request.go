package letsrest

// информация задаче на выполнение запроса
type RequestTask struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

// информация об ответе на запрос
type Response struct {
	ID string `json:"id"`

	StatusCode int `json:"status_code"`
}

// информация о статусе выполнения запроса вместе с ответом
type Result struct {
	Status   *Info     `json:"info"`
	Response *Response `json:"response,omitempty"`
}

// информация о статусе запроса
type Info struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
