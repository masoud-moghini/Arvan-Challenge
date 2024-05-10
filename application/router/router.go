package router

import (
	"net/http"

	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	Router interface {
		Listen()
	}
	RouterHandler struct {
		routes          *chi.Mux
		RequestHandlers RequestHandlers
	}
)

func (routerHandler RouterHandler) Listen() {
	routerHandler.routes = chi.NewRouter()
	routerHandler.routes.Use(middleware.Logger)
	routerHandler.routes.Post("/", func(w http.ResponseWriter, r *http.Request) {
		routerHandler.RequestHandlers.HandleIncommingRequest(w, r)

	})
	http.ListenAndServe(":3000", routerHandler.routes)
	fmt.Println("hi")
}
