package model

import "github.com/vkosev/ft_api/internal/expression"

type ExpressionRequest struct {
	// Text expression to be resolved
	// example: What is 5 plus 3 multiplied by 2?
	// required: true
	// in:body
	Expression string `json:"expression" validate:"required"`
}

// swagger:model EvaluatedExpression
type EvaluatedExpressionResponse struct {
	// Result of the evaluated expression
	// example: 11
	Result int `json:"result"`
}

// swagger:model ExpressionResult
type InvalidExpressionResponse struct {
	// Boolean value true if valid, false if not
	// example: false
	Valid bool `json:"valid"`
	// Reason why expression is not valid
	// required: false
	Reason string `json:"reason"`
}

// swagger:model ValidExpression
type ValidExpressionResponse struct {
	// Boolean value true if valid, false if not
	Valid bool `json:"valid"`
}

// swagger:model ExpressionError
type ExpressionErrorResponse struct {
	// Type of the error
	Type expression.ErrorType `json:"type"`
	// Message of the error
	Message string `json:"message"`
}

// swagger:model Error
type ErrorHistoryResponse struct {
	// The expression for which the error occurred
	Expression string `json:"expression"`
	// The endpoints on which the error occurred
	// collection format: slice
	Endpoints []ErrorEndpoint `json:"endpoints"`
	// Count of the number of times the error occurred for this expression
	Frequency int `json:"frequency"`
	// Type of the error
	ErrType expression.ErrorType `json:"type"`
}

type ErrorEndpoint struct {
	Url   string
	Count int
}
