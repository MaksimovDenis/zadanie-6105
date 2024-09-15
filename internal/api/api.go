package api

import (
	"context"
	"net/http"
	"time"

	"git.codenrock.com/zadanie-6105/internal/storage"
	"git.codenrock.com/zadanie-6105/pkg/protocol/oapi"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	FileUploadBufferSize       = 512e+6 // 512MB for now
	ServerShutdownDefaultDelay = 5 * time.Second
)

type Opts struct {
	Addr    string
	Log     zerolog.Logger
	Storage *storage.Storage
}

type API struct {
	l       zerolog.Logger
	server  *http.Server
	router  *gin.Engine
	storage *storage.Storage
}

func NewAPI(opts *Opts) (*API, error) {
	router := gin.Default()

	router.MaxMultipartMemory = FileUploadBufferSize

	api := &API{
		l: opts.Log,
		server: &http.Server{
			Addr:    opts.Addr,
			Handler: router,
		},
		router:  router,
		storage: opts.Storage,
	}

	router.Use(TokenInjector())

	oapi.RegisterHandlersWithOptions(router, api, oapi.GinServerOptions{
		BaseURL: "/api",
	})

	return api, nil
}

func (api *API) Serve() error {
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to start api server")
		return err
	}

	return nil
}

func (api *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownDefaultDelay)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		api.l.Error().Err(err).Msg("failed to stop api server")
	}
}
