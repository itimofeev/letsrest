package letsrest

type ClientRequest struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Method string `json:"method"`
}

type ClientResponse struct {
	StatusCode int `json:"status_code"`
}
