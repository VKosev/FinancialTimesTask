package routes

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/vkosev/ft_api/internal/expression"
	"github.com/vkosev/ft_api/internal/web/handler"
)

func Init(logger *log.Logger, resolver *expression.Resolver) *mux.Router {
	h := handler.NewHandler(logger, resolver)

	router := mux.NewRouter()

	router.HandleFunc("/expression", h.Evaluate).Methods("POST")
	router.HandleFunc("/errors", h.AllErrors).Methods("GET")
	// router.HandleFunc("/validate", h.).Methods("POST")

	return router
}
