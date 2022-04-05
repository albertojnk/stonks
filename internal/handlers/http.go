package handlers

import (
	"net/http"

	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/ports"
	"github.com/gin-gonic/gin"
)

//HTTPHandler .
type HTTPHandler struct {
	authService  ports.AuthService
	stockService ports.StockService
}

//NewHTTPHandler new account handler
func NewHTTPHandler(
	authService ports.AuthService,
	stockService ports.StockService,
) *HTTPHandler {
	return &HTTPHandler{
		authService:  authService,
		stockService: stockService,
	}
}

// LoginHandler represent endpoint: [GET] /dashboards
func (hdl *HTTPHandler) LoginHandler(ctx *context.Context, c *gin.Context) error {
	body := gin.H{}

	body["http_prefix"] = ctx.HTTPPrefix

	c.HTML(http.StatusOK, "login.html", body)
	return nil
}
