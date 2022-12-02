package service

// Evaluator responsible to calculate arithmetic expression from text
type Evaluator interface {
	Evaluate(expr string) (int error)
}

type ExpressionResolver struct {
	evaluator Evaluator
}

func NewExpressionResolver(evaluator Evaluator) *ExpressionResolver {
	return &ExpressionResolver{
		evaluator: evaluator,
	}
}

func (er *ExpressionResolver) Evaluate(expr string) (int, error) {
	return 11, nil
}
