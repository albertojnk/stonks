package cmd

import (
	systemcontext "context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/albertojnk/stonks/internal/common"
	"github.com/albertojnk/stonks/internal/context"
	"github.com/albertojnk/stonks/internal/core/services/authservice"
	"github.com/albertojnk/stonks/internal/core/services/stockservice"
	"github.com/albertojnk/stonks/internal/db"
	"github.com/albertojnk/stonks/internal/handlers"
	"github.com/albertojnk/stonks/internal/repositories/authrepository"
	"github.com/albertojnk/stonks/internal/repositories/stockrepository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/unrolled/secure"
)

var api = &cobra.Command{
	Use:   "web",
	Short: "Start web application",
	RunE: func(cmd *cobra.Command, args []string) error {

		systemctx, cancel := systemcontext.WithCancel(systemcontext.Background())
		// --> Context
		ctx := context.New().WithLogger()
		ctx.HTTPPrefix = common.GetEnv("HTTPPREFIX", "")

		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			ctx.Logger.Info("Shutting down API, waiting...")
			cancel()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveWeb(ctx, systemctx)
		}()

		wg.Wait()
		return nil
	},
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func serveWeb(ctx *context.Context, systemctx systemcontext.Context) error {

	// --> Database
	dbConfig := db.NewPostgresConfig()
	db, err := db.NewPostgres(dbConfig)
	if err != nil {
		ctx.Logger.Errorf("Error initilize database instance - Error: %v", err.Error())
		panic(err)
	}

	// --> Auth
	authrepo := authrepository.New(db)
	authsrv := authservice.New(authrepo)

	// --> Stock

	stockrepo := stockrepository.New(db)
	stocksrv := stockservice.New(stockrepo)

	// --> Handler
	hdl := handlers.NewHTTPHandler(authsrv, stocksrv)

	router := gin.New()
	env := common.GetEnv("ENVIRONMENT", "dev")
	if env == "prod" {
		router.Use(func() gin.HandlerFunc {
			return func(c *gin.Context) {
				secureMiddleware := secure.New(secure.Options{
					SSLRedirect: true,
					SSLHost:     "", // url
				})
				err := secureMiddleware.Process(c.Writer, c.Request)

				if err != nil {
					return
				}

				c.Next()
			}
		}())
	}

	// router.Use(sessions.Sessions(domains.DefaultSessionID, store))
	router.Use(handlers.GlobalMiddleware())
	router.Static(ctx.HTTPPrefix+"/assets", "web/assets")
	router.SetFuncMap(template.FuncMap{
		"dict": dict,
	})
	router.LoadHTMLGlob("web/pages/*")

	// --> API's
	router.GET(ctx.HTTPPrefix+"/api/v1/auth", handlers.HandlerAPI(true, hdl.AuthLogin))
	router.GET(ctx.HTTPPrefix+"/api/v1/stock", handlers.HandlerAPI(false, hdl.StockGet))

	// --> PAGES
	router.GET(ctx.HTTPPrefix+"/login", handlers.HandlerPage(false, hdl.LoginHandler))

	port := common.GetEnv("PORT", "9000")
	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	done := make(chan struct{})
	doneTLS := make(chan struct{})
	go func() {
		<-systemctx.Done()
		if err := s.Shutdown(systemcontext.Background()); err != nil {
			ctx.Logger.Error(err)
		}
		close(done)
	}()

	if env == "prod" {
		go func() {
			ctx.Logger.Infof("Serving web application at http://0.0.0.0:%d", 80)
			if err := s.ListenAndServe(); err != http.ErrServerClosed {
				ctx.Logger.Error(err)
			}
		}()

		sTLS := &http.Server{
			Addr:    ":443",
			Handler: router,
		}

		go func() {
			<-systemctx.Done()
			if err := sTLS.Shutdown(systemcontext.Background()); err != nil {
				ctx.Logger.Error(err)
			}
			close(doneTLS)
		}()

		ctx.Logger.Infof("Serving web application at http://0.0.0.0:%s", "443")
		if err := sTLS.ListenAndServeTLS("./configs/tls.crt", "./configs/tls.key"); err != http.ErrServerClosed {
			ctx.Logger.Error(err)
		}
	} else {
		close(doneTLS)
		ctx.Logger.Infof("Serving web application at http://0.0.0.0:%v", port)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			ctx.Logger.Error(err)
		}
	}

	<-doneTLS
	<-done

	return nil
}

func init() {
	rootCmd.AddCommand(api)
}
