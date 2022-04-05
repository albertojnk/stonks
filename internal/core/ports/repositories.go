package ports

import (
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
)

type AuthRepository interface {
	Login(*context.Context, *domains.Auth) (context.Result, *domains.Auth)
}

type StockRepository interface {
	Get(*context.Context, string, string) (context.Result, domains.MarketQuotesResponse)
}
