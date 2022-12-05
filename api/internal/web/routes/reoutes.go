package routes

import (
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/vkosev/ft_api/internal/expression"
	"github.com/vkosev/ft_api/internal/web/handler"
)

func Init(logger *log.Logger, resolver *expression.Resolver) *mux.Router {
	h := handler.NewHandler(logger, resolver)

	router := mux.NewRouter()

	// swagger:route POST /expression Expression
	//
	// Evaluates an arithmetic text expression
	//
	//
	// responses:
	//   200: EvaluatedExpression
	//   400: ExpressionError
	//   500: Message
	router.HandleFunc("/expression", h.Evaluate).Methods("POST")

	// swagger:route POST /validate ValidExpression
	//
	// Checks wether expression is valid
	//
	//
	// responses:
	//   200: ValidExpression
	//   400: InvalidExpression
	//   500: Message
	router.HandleFunc("/validate", h.Validate).Methods("POST")

	// swagger:route GET /errors ErrorHistory
	//
	// Returns all occurred expression errors
	//
	//
	// responses:
	//   200: ExpressionErrors
	//   500: Message
	router.HandleFunc("/errors", h.AllErrors).Methods("GET")

	registerSwaggerEndpoint(router)

	return router
}

func registerSwaggerEndpoint(router *mux.Router) {
	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("../../")))
}
