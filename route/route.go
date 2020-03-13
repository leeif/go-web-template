package route

import (
	"net/http"

	"github.com/leeif/go-web-template/config"
	"github.com/leeif/go-web-template/datatype/error"
	"github.com/leeif/go-web-template/middleware"

	"github.com/gorilla/mux"
	"github.com/leeif/go-web-template/log"
	"github.com/leeif/go-web-template/manage"

	"github.com/urfave/negroni"
)

type middle func(handlers ...negroni.HandlerFunc) http.Handler

type RouteItem struct {
	path        string
	description string
	method      string
	middle      middle
	handler     func(w http.ResponseWriter, r *http.Request) *error.Error
}

type Route struct {
	manager *manage.Manager
	config  *config.Config
	logger  *log.Log
	mw      *middleware.Middleware
	mux     *mux.Router
}

func (route *Route) registerHealthRoutes() {
	routes := []RouteItem{
		{
			path:        "/healthcheck",
			description: "health check api",
			method:      "GET",
			middle:      route.mw.NoAuthMiddleware,
			handler:     route.healthcheck,
		},
	}
	route.registerRoutes(routes, "/", false)
}

func (route *Route) registerRoutes(routes []RouteItem, prefix string, isWeb bool) {
	sub := route.mux.PathPrefix(prefix).Subrouter()
	for _, r := range routes {
		// options method for cors
		if isWeb {
			sub.Handle(r.path, r.middle(route.handlerWrapper(r.handler))).Methods(r.method)
		} else {
			sub.Handle(r.path, r.middle(route.webHandlerWrapper(r.handler))).Methods(r.method)
		}
	}
}

func (route *Route) Register() {
	route.registerHealthRoutes()
}

func NewRouter(mux *mux.Router, manager *manage.Manager, config *config.Config, logger *log.Log) *Route {
	return &Route{
		manager: manager,
		config:  config,
		logger:  logger,
		mw:      middleware.NewMiddle(logger, config),
		mux:     mux,
	}
}
