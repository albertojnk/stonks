package stockrepository

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

func (r *Repository) Get(ctx *context.Context, stock string, region string) (context.Result, domains.MarketQuotesResponse) {
	return ctx.ResultSuccess(), domains.MarketQuotesResponse{}
}
