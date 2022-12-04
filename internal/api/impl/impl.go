package impl

import (
	"github.com/apachejuice/eelchat/internal/api"
	"github.com/apachejuice/eelchat/internal/db/repository"
	"github.com/apachejuice/eelchat/internal/logs"
)

type Impl struct {
	userRepo repository.UserRepository
	log      logs.Logger
}

func NewImpl() Impl {
	return Impl{
		userRepo: repository.NewUserRepository(),
		log:      logs.NewLogger(api.ComponentImpl),
	}
}
