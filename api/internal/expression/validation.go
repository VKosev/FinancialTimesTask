package expression

import "regexp"

func isValid(expr string) error {
	if err := isNonMathQuestion(expr); err != nil {
		return err
	}

	err := isUnsupportedOperations(expr)
	if err != nil {
		return err
	}

	return nil
}

func isNonMathQuestion(expr string) error {
	startWithMathQuestionPattern := regexp.MustCompile(`What\sis\s\d+\s(plus|multiplied|minus|divided)\s\d+`)
	matchNumbersPattern := regexp.MustCompile(matchNumbersRegexExpression)

	if !startWithMathQuestionPattern.MatchString(expr) && !matchNumbersPattern.MatchString(expr) {
		return ExpressionError{
			Msg:        "Expression is not arithmetic question and does not contain numbers",
			Expression: expr,
			ErrType:    NonMathQuestion,
		}
	}

	return nil
}

func isUnsupportedOperations(expr string) error {
	startAsExpressionPattern := regexp.MustCompile(`What\sis\s\d+\s`)
	endsWithNumberAndQuestionMarkPattern := regexp.MustCompile(`\d+\?$|\w+\?$`)

	if startAsExpressionPattern.MatchString(expr) && endsWithNumberAndQuestionMarkPattern.MatchString(expr) {
		containsSupportedOperationsPattern := regexp.MustCompile(`plus\s\d+|minus\s\d+|multiplied\sby\s\d+|divided\sby\s\d+`)

		if !containsSupportedOperationsPattern.MatchString(expr) {
			return ExpressionError{
				Msg:        "Expression contains unsupported operations",
				Expression: expr,
				ErrType:    UnsupportedOperations,
			}
		}
	}

	return nil
}
