package authrepository

import (
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
	"github.com/albertojnk/stonks/internal/core/ports"
)

type Repository struct {
	Db ports.Persistance
}

func New(db ports.Persistance) *Repository {
	return &Repository{
		Db: db,
	}
}

func (r *Repository) Login(ctx *context.Context, auth *domains.Auth) (context.Result, *domains.Auth) {
	return ctx.ResultSuccess(), auth
}
