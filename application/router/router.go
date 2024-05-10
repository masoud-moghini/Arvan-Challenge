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
		Routes          *chi.Mux
		RequestHandlers RequestHandlers
	}
)

func (routerHandler RouterHandler) Listen() {
	routerHandler.Routes.Use(middleware.Logger)
	routerHandler.Routes.Post("/", func(w http.ResponseWriter, r *http.Request) {
		routerHandler.RequestHandlers.HandleIncommingRequest(w, r)

	})
	http.ListenAndServe(":3000", routerHandler.Routes)
	fmt.Println("hi")
}
