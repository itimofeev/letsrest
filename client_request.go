package letsrest

type ClientRequest struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

type ClientResponse struct {
	ID string `json:"id"`

	StatusCode int `json:"status_code"`
}

type RequestData struct {
	ID string `json:"id"`

	Request  *ClientRequest
	Response *ClientResponse
}
