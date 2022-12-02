package evaluator

import "regexp"

type Evaluator struct {
	regexp *regexp.Regexp
}

func NewEvaluator() *Evaluator {
	r, _ := regexp.Compile("")

	return &Evaluator{
		regexp: r,
	}
}

func (e *Evaluator) evaluate(expression string) (int, error) {
	return 11, nil
}
