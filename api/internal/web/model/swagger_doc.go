// swagger_doc.go is used only for generating
// swagger spec based on the models here.
// The models here are not used anywhere else

package model

// swagger:parameters Expression
type RequestBody struct {
	// in: body
	Body ExpressionRequest
}

// swagger:model ExpressionErrors
type ErrorHistoriesResponse struct {
	// collection format: slice
	Errors []ErrorHistoryResponse
}

// swagger:model Message
type _ string
