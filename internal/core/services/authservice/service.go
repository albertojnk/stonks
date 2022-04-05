package authservice

import (
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
	"github.com/albertojnk/stonks/internal/core/ports"
)

type Service struct {
	authService ports.AuthRepository
}

func New(repository ports.AuthRepository) *Service {
	return &Service{
		authService: repository,
	}
}

func (s *Service) Login(ctx *context.Context, auth *domains.Auth) (context.Result, *domains.Auth) {
	return s.authService.Login(ctx, auth)
}
