package rest

import "github.com/apachejuice/eelchat/internal/api/spec"

type baseCtx struct {
	actionId ActionID
}

func (baseCtx) Error(msg, component string) spec.Error {
	return spec.Error{
		Message:   msg,
		Component: component,
	}
}

func (baseCtx) FromError(err error, component string) spec.Error {
	return spec.Error{
		Message:   "Internal server error: " + err.Error(),
		Component: component,
	}
}

// A type that helps generating CreateUser responses.
type CreateUserContext struct{ baseCtx }

func (CreateUserContext) NoContent() *spec.CreateUserNoContent {
	return &spec.CreateUserNoContent{}
}

func (c CreateUserContext) BadRequest(e spec.Error) *spec.CreateUserApplicationJSONBadRequest {
	e.ActionId = c.actionId.String()
	return (*spec.CreateUserApplicationJSONBadRequest)(&e)
}

func (c CreateUserContext) InternalServerError(e spec.Error) *spec.CreateUserApplicationJSONInternalServerError {
	e.ActionId = c.actionId.String()
	return (*spec.CreateUserApplicationJSONInternalServerError)(&e)
}
