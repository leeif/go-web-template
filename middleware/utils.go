package middleware

import (
	"net/http"

	"github.com/leeif/go-web-template/config"
	"github.com/leeif/go-web-template/log"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger *log.Log
	Config *config.Config
}

func (middleware *Middleware) NoAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	return ng
}

func NewMiddle(logger *log.Log, config *config.Config) *Middleware {
	return &Middleware{
		Logger: logger.With("componment", "middleware"),
		Config: config,
	}
}
