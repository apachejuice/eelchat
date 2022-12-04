package impl

import (
	"github.com/apachejuice/eelchat/internal/api"
	"github.com/apachejuice/eelchat/internal/api/crypt"
	"github.com/apachejuice/eelchat/internal/api/rest"
	"github.com/apachejuice/eelchat/internal/api/spec"
)

// Implements creating users.
func (i Impl) CreateUserImpl(ctx rest.CreateUserContext, rctx rest.RequestContext, reqBody spec.User) spec.CreateUserRes {
	hash, err := crypt.Hash(reqBody.Password)
	if err != nil {
		return ctx.InternalServerError(ctx.FromError(err, api.ComponentImpl))
	}

	if hash == "" || len(reqBody.Username) < 3 {
		return ctx.BadRequest(ctx.Error("Username must be at least 3 characters long and password at least 8 characters long", api.ComponentImpl))
	}

	user, err := i.userRepo.Create(
		reqBody.Username,
		hash,
		reqBody.Email.Or(""),
	)

	if err != nil {
		return ctx.InternalServerError(ctx.FromError(err, api.ComponentImpl))
	}

	i.log.Debug("User created", "username", user.Username, "discriminator", user.Discriminator)
	return ctx.NoContent()
}
