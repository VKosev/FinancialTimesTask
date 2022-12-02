package model

type ExpressionRequest struct {
	expression string `json:"expression"`
}

type EvaluatedExpressionResponse struct {
	result int `json:"result"`
}

type InvalidExpressionResponse struct {
	valid  bool   `json:"valid"`
	reason string `json:"reason"`
}

type ValidExpressionResponse struct {
	valid bool `json:"valid"`
}

type ErrorResponse struct {
	expression string `json:"expression"`
	endpoint   string `json:"endpoint"`
	frequency  int    `json:"frequency"`
	errType    string `json:"type"`
}
