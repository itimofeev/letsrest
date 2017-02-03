package letsrest

type ClientRequest struct {
	URL    string `json:"url" binding:"required"`
	Method string `json:"method"`
}

type ClientResponse struct {
	StatusCode int `json:"status_code"`
}
