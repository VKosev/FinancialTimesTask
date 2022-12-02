package handler

import (
	"log"
	"net/http"

	"github.com/vkosev/ft_api/internal/evaluator"
)

type Handler struct {
	l        *log.Logger
	resolver *evaluator.Evaluator
}

func (h *Handler) NewHandler(l *log.Logger, resolver *evaluator.Evaluator) *Handler {
	return &Handler{
		l:        l,
		resolver: resolver,
	}
}

func (h *Handler) evaluate(w http.ResponseWriter, req *http.Request) {

}
