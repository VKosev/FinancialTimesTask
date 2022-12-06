package expression

type ExpressionRequest struct {
	Expression string `json:"expression" validate:"required"`
}

type EvaluatedExpressionResponse struct {
	Result int `json:"result"`
}

type InvalidExpressionResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

type ValidationResult struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

type ValidExpressionResponse struct {
	Valid bool `json:"valid"`
}

type EvaluationResult struct {
	Result  int
	Type    string
	Message string
}

type FailEvaluationError string

func (e FailEvaluationError) Error() string {
	return "Evaluation failed."
}

type ExpressionErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type ErrorHistoryResponse struct {
	Expression string          `json:"expression"`
	Endpoints  []ErrorEndpoint `json:"endpoints"`
	Frequency  int             `json:"frequency"`
	ErrType    string          `json:"errorType"`
}

type ErrorEndpoint struct {
	Url   string
	Count int
}
