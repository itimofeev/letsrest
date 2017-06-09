package letsrest

type ResponseData struct {
	Response *Response   `json:"response,omitempty"`
	Status   *ExecStatus `json:"status"`
}
