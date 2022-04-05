package handlers

import (
	"net/http"
	"time"

	"github.com/albertojnk/stonks/internal/common"
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/domains"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (hdl *HTTPHandler) AuthLogin(ctx *context.Context, c *gin.Context) error {
	params := &domains.Auth{}
	err := c.ShouldBindJSON(params)
	if err != nil {
		ctx.ResultError(1, "Invalid Request").JSON(c, nil)
		return nil
	}

	result, auth := hdl.authService.Login(ctx, params)
	if result.State != context.ResultStateSuccess {
		ctx.ResultError(0, "error to authorize").JSON(c, nil)
	}

	claims := &domains.JWTClaims{
		*auth,
		false,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(common.GetEnv("SECRET_TOKEN", "word_battle_123")))
	if err != nil {
		return err
	}

	session := sessions.Default(c)
	session.Set("access_token", t)
	session.Delete("state")
	err = session.Save()
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": t,
		"token_type":   "Bearer",
		"expires_in":   claims.ExpiresAt,
	})

	return nil
}
