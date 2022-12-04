package rest

import (
	"context"

	"github.com/apachejuice/eelchat/internal/api/spec"
)

type (
	RequestContext = context.Context // clearly specifies the context's purpose

	CreateUserFunc func(CreateUserContext, RequestContext, spec.User) spec.CreateUserRes
)

// Helper type to connect the spec
type apiHandler struct {
	// Creates a new user.
	createUserFunc CreateUserFunc

	lastActionId ActionID // Assigned in api.go
}

var _ spec.Handler = (*apiHandler)(nil)

func (a apiHandler) CreateUser(reqCtx context.Context, user spec.User) (spec.CreateUserRes, error) {
	res := a.createUserFunc(CreateUserContext{baseCtx{a.lastActionId}}, reqCtx, user)
	return res, nil
}
