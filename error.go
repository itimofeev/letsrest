package letsrest

import "fmt"

const ReqNotFoundKey = "request.not.found"

func RequestNotFoundResponse(id string) *ErrorResponse {
	return &ErrorResponse{
		Key:         ReqNotFoundKey,
		Description: fmt.Sprintf("Request with id %s not found in store", id),
		Params:      Params{"id": id},
	}
}

type Params map[string]string

type ErrorResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Params      Params `json:"params"`
}
