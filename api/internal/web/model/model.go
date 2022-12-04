package model

import "github.com/vkosev/ft_api/internal/expression"

// swagger:model ExpressionRequest
type ExpressionRequest struct {
	// Text expression to be resolved
	Expression string `json:"expression"`
}

// swagger:model EvaluatedExpressionResponse
type EvaluatedExpressionResponse struct {
	// Result of the evaluated expression
	Result int `json:"result"`
}

type InvalidExpressionResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

type ValidExpressionResponse struct {
	Valid bool `json:"valid"`
}

type ExpressionErrorResponse struct {
	Type    expression.ErrorType `json:"type"`
	Message string               `json:"message"`
}

type ErrorHistoryResponse struct {
	Expression string               `json:"expression"`
	Endpoints  []ErrorEndpoint      `json:"endpoints"`
	Frequency  int                  `json:"frequency"`
	ErrType    expression.ErrorType `json:"type"`
}

type ErrorEndpoint struct {
	Url   string
	Count int
}
