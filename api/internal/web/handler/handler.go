package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator"
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

func (h *Handler) Evaluate(w http.ResponseWriter, req *http.Request) {
	reqBody := &model.ExpressionRequest{}
	if ok := decodeAndValidate(w, req, reqBody); !ok {
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

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, model.EvaluatedExpressionResponse{
		Result: evaluatedExpression,
	})
}

func (h *Handler) Validate(w http.ResponseWriter, req *http.Request) {
	reqBody := &model.ExpressionRequest{}
	if ok := decodeAndValidate(w, req, reqBody); !ok {
		return
	}

	err := h.resolver.Validate(reqBody.Expression, req.URL.Path)
	if err != nil {
		if e, ok := expression.IsExpressionError(err); ok {
			writeJSON(w, http.StatusOK, model.InvalidExpressionResponse{
				Valid:  false,
				Reason: e.Msg,
			})
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func decodeAndValidate(w http.ResponseWriter, req *http.Request, v any) bool {
	err := json.NewDecoder(req.Body).Decode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	validate := validator.New()
	if err = validate.Struct(v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
