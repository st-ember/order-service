package api

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/st-ember/ecommerceprocessor/internal/processor/handler"
)

func Router(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/purchase", func(router chi.Router) {
		router.Post("/", handler.CreatePurchase)
		router.Get("{id}/status", handler.PurchaseStatus)
		router.Delete("/{id}", handler.DeletePurchase)
	})
}
