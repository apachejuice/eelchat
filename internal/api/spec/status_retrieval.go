package spec

import "net/http"

// This file is in the spec package to implement retrieving the status code of a response.

func (*CreateUserNoContent) GetStatusCode() int                 { return http.StatusNoContent }
func (*CreateUserApplicationJSONBadRequest) GetStatusCode() int { return http.StatusBadRequest }
func (*CreateUserApplicationJSONInternalServerError) GetStatusCode() int {
	return http.StatusInternalServerError
}
