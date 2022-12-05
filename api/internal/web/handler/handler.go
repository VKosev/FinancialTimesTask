package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vkosev/ft_api/internal/expression"
	"github.com/vkosev/ft_api/internal/web/model"
)

type Handler struct {
	l        *log.Logger
	resolver *expression.Resolver
}

// NewHandler returns a pointer to new Handler instance
func NewHandler(l *log.Logger, resolver *expression.Resolver) *Handler {
	return &Handler{
		l:        l,
		resolver: resolver,
	}
}

// swagger:route POST /posts PostsList
// Evaluates an arithmetic expression
//
// responses:
//  200: PostResponse
func (h *Handler) Evaluate(w http.ResponseWriter, req *http.Request) {
	reqBody := &model.ExpressionRequest{}

	err := json.NewDecoder(req.Body).Decode(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	evaluatedExpression, err := h.resolver.Evaluate(reqBody.Expression, req.URL.Path)
	if err != nil {
		if e, ok := expression.IsExpressionError(err); ok {
			writeJSON(w, http.StatusBadRequest, model.ExpressionErrorResponse{
				Type:    e.ErrType,
				Message: e.Error(),
			})
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, model.EvaluatedExpressionResponse{
		Result: evaluatedExpression,
	})
}

func (h *Handler) Validate(w http.ResponseWriter, req *http.Request) {
	reqBody := &model.ExpressionRequest{}

	err := json.NewDecoder(req.Body).Decode(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.resolver.Validate(reqBody.Expression, req.URL.Path)
	if err != nil {
		if e, ok := expression.IsExpressionError(err); ok {
			writeJSON(w, http.StatusBadRequest, model.InvalidExpressionResponse{
				Valid:  false,
				Reason: e.Msg,
			})

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, model.ValidExpressionResponse{
		Valid: true,
	})
}

func (h *Handler) AllErrors(w http.ResponseWriter, req *http.Request) {
	errors := h.resolver.ErrorHistory()

	responseErrors := make([]model.ErrorHistoryResponse, len(errors))

	for i, err := range errors {
		responseErrors[i] = model.ErrorHistoryResponse{
			Expression: err.Expression,
			Endpoints:  convertToErrorEndpoints(err.Endpoints),
			Frequency:  err.Frequency,
			ErrType:    err.ErrType,
		}
	}

	writeJSON(w, http.StatusOK, responseErrors)
}

func convertToErrorEndpoints(m map[string]int) []model.ErrorEndpoint {
	endpoints := make([]model.ErrorEndpoint, len(m))

	index := 0
	for url, count := range m {
		errEndpoint := model.ErrorEndpoint{
			Url:   url,
			Count: count,
		}
		endpoints[index] = errEndpoint

		index++
	}

	return endpoints
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Marshalling JSON failed.", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
