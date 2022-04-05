package ports

import (
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
)

type AuthService interface {
	Login(*context.Context, *domains.Auth) (context.Result, *domains.Auth)
}

type StockService interface {
	Get(*context.Context, string, string) (context.Result, domains.MarketQuotesResponse)
}
