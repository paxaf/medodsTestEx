package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/paxaf/medodsTestEx/config"
	_ "github.com/paxaf/medodsTestEx/docs"
	"github.com/paxaf/medodsTestEx/internal/controller/httpserver"
	"github.com/paxaf/medodsTestEx/internal/repository"
	"github.com/paxaf/medodsTestEx/internal/usecase"
)

type App struct {
	config    *config.Config
	apiServer *http.Server
}

func New(cfg *config.Config) (*App, error) {
	app := &App{}
	app.config = cfg

	pool, err := pgxpool.New(context.Background(), cfg.Database.GetDSN())
	if err != nil {
		log.Error().Err(err)
	}
	repo := repository.NewRepository(pool)

	usecase := usecase.NewUseCase(repo)
	handler := httpserver.NewFeedbackHandler(usecase)
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", handler.SubmitPing)
	router.GET("/tokens", handler.GetTokens)
	router.GET("/guid", handler.Guid)
	router.POST("/refresh", handler.Refresh)
	router.GET("/logout", handler.Deauthorized)
	host := app.config.APIServer.Host
	port := app.config.APIServer.Port
	addr := net.JoinHostPort(host, port)
	app.apiServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	log.Info().Msg("app initialized successfully")
	return app, nil
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Msg("API server started successfuly " + "address " + app.apiServer.Addr)
		if err := app.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err)
		}
	}()

	<-ctx.Done()
	log.Info().Msg("received shutdown signal")
	return nil
}
