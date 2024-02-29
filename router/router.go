package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"playlistturbo.com/controller"
	"playlistturbo.com/router/middlewares"
)

func Get(ctrl controller.Controller) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)

	applyRoutes(r, ctrl, "/songs", SongsRoutes(ctrl))

	return r
}

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	// Auth    *middlewares.Auth
	Params middlewares.Params
	Body   interface{}
}

func applyRoutes(r *chi.Mux, ctrl controller.Controller, subpath string, endpoints []Route) {
	subrouter := r
	// create a subrouter only if subpath is not empty
	if subpath != "" {
		subrouter = chi.NewRouter()
	}

	for _, e := range endpoints {
		handler := e.Handler

		// apply middlewares in reverse order because they're applied from the outer to the inner one

		if len(e.Params) > 0 || e.Body != nil {
			m := middlewares.ParamsMiddleware(e.Params, e.Body)
			handler = m(handler).ServeHTTP // apply middleware
		}

		// apply logger middleware
		// m := middlewares.LoggerMiddleware()
		// handler = m(handler).ServeHTTP // apply middleware

		// // TO FIX use AuthMiddleware no checkAuthMdlw
		// if e.Auth != nil {
		// 	m := middlewares.AuthMiddleware(*e.Auth)
		// 	handler = m(handler).ServeHTTP // apply middleware

		// 	// apply CreateOrUpdate user
		// 	m = middlewares.CheckCreateUserOrUpdateMiddle(ctrl)
		// 	handler = m(handler).ServeHTTP // apply middleware
		// }

		// // apply panic middleware
		// m = middlewares.PanicMiddleware()
		// handler = m(handler).ServeHTTP // apply middleware

		// // apply ID middleware
		// m = middlewares.IDMiddleware()
		// handler = m(handler).ServeHTTP // apply middleware

		// set method, path, and handler for this endpoint
		subrouter.Method(e.Method, e.Path, handler)
	}

	// mount a subrouter only if subpath is not empty
	if subpath != "" {
		r.Mount(subpath, subrouter)
	}
}
