package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/albertojnk/stonks/internal/common"
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/i18n"
	"github.com/gin-gonic/gin"
)

const RequestContextKey = "request_context"

// HandlerPage provide abstractionn for page handler receive the context of application.
func HandlerPage(mustBeAuthenticate bool, f func(ctx *context.Context, c *gin.Context) error) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctxVal := c.MustGet(RequestContextKey).(*context.Context)

		if mustBeAuthenticate && ctxVal.LoggedAuth.ID == 0 {
			c.Redirect(303, ctxVal.HTTPPrefix+"/login")
			return
		}

		if err := f(ctxVal, c); err != nil {
			ctxVal.Logger.Error(fmt.Errorf("%v: %s", err, debug.Stack()))
			c.Redirect(302, "/error")
		}
	})
}

//HandlerAPI .
func HandlerAPI(mustBeAuthenticate bool, f func(ctx *context.Context, c *gin.Context) error) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctxVal := c.MustGet(RequestContextKey).(*context.Context)

		if mustBeAuthenticate && ctxVal.LoggedAuth.ID == 0 {
			c.Status(401)
			return
		}

		if err := f(ctxVal, c); err != nil {
			ctxVal.Logger.Error(fmt.Errorf("%v: %s", err, debug.Stack()))
			http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
		}
	})
}

func GlobalMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// lang := findLanguage(c)

		ctx := context.New().WithLogger()
		ctx.HTTPPrefix = common.GetEnv("HTTPPREFIX", "")
		ctx.HostURL = c.Request.Host
		// ctx.Lang = context.LangType(lang)
		// i18n.SetLang(ctx, lang)

		defer func() {
			if r := recover(); r != nil {
				http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		c.Set(RequestContextKey, ctx)
		c.Next()
	})
}

func findLanguage(c *gin.Context) string {
	lang := ""
	langKeys := common.GetEnv("LANGKEYS", "lang,language")
	langKeysSlice := strings.Split(langKeys, ",")

	for _, key := range langKeysSlice {
		var hasQuery, hasCookie bool
		langCookie, _ := c.Cookie(key)
		langQuery := c.Query(key)

		if langCookie != "" {
			hasCookie = true
		}
		if langQuery != "" {
			hasQuery = true
		}

		switch {
		case hasQuery && hasCookie:
			lang = i18n.GetSupportedLang(langQuery)
		case hasQuery && !hasCookie:
			lang = i18n.GetSupportedLang(langQuery)
		case !hasQuery && hasCookie:
			lang = i18n.GetSupportedLang(langCookie)
		default:
			continue
		}

		if lang != "" {
			c.SetCookie(key, lang, 20000, "/", c.Request.URL.Host, false, false)
			break
		}
	}
	if lang == "" {
		lang = "en-US"
		c.SetCookie(langKeysSlice[0], lang, 20000, "/", c.Request.URL.Host, false, false)
	}

	return i18n.GetSupportedLang(lang)
}
