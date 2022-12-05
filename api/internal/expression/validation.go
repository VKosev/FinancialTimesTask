package expression

import "regexp"

var startAsExpressionPattern *regexp.Regexp = regexp.MustCompile(`What\sis\s\d+\s`)
var endsWithNumberAndQuestionMarkPattern *regexp.Regexp = regexp.MustCompile(`\d+\?$|\w+\?$`)

// isValid validates wether the passes string is valid expressions text.
//
// Returns nil if valid, otherwise error which can be casted to ExpressionError.
func isValid(expr string) error {
	if err := isNonMathQuestion(expr); err != nil {
		return err
	}

	if err := isUnsupportedOperations(expr); err != nil {
		return err
	}

	if err := isInvalidSyntaxError(expr); err != nil {
		return err
	}

	return nil
}

// isNonMathQuestion checks if the passed string is non mathematical questions.
//
// Returns error if true, otherwise nil
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

// isUnsupportedOperations check if the passed string contains unsupported operations.
//
// Returns error if true, otherwise nil
func isUnsupportedOperations(expr string) error {
	containsSupportedOperationsPattern := regexp.MustCompile(`plus\s\d+|minus\s\d+|multiplied\sby\s\d+|divided\sby\s\d+`)

	if startAsExpressionPattern.MatchString(expr) && endsWithNumberAndQuestionMarkPattern.MatchString(expr) {
		if !containsSupportedOperationsPattern.MatchString(expr) {
			return NewUnsupportedExpressionError(expr)
		}
	}

	finalOperationIsSupported := regexp.MustCompile(`multiplied\sby\s\d+\?$|divided\sby\s\d+\?$|plus\s\d+\?$|minus\s\d+\?$`)
	if containsSupportedOperationsPattern.MatchString(expr) && !finalOperationIsSupported.MatchString(expr) {
		return NewUnsupportedExpressionError(expr)
	}

	return nil
}

func isInvalidSyntaxError(expr string) error {
	if !isNonMathQuestionOrUnsupportedOperation(expr) {
		validAndInvalidOperations := regexp.MustCompile(`plus[a-zA-Z a-zA-Z]*|minus[a-zA-Z a-zA-Z]*|multiplied\sby[a-zA-Z a-zA-Z]*|divided\sby[a-zA-Z a-zA-Z]*`)

		matchedStrings := validAndInvalidOperations.FindAllString(expr, -1)

		if len(matchedStrings) == 0 {
			return NewInvalidSyntaxExpressionError(expr, "Expression does not contain operations")
		}

		if !containsOnlySupportedOperations(matchedStrings) {
			return NewInvalidSyntaxExpressionError(expr, "Expression is not valid")
		}

	}

	return nil
}

func containsOnlySupportedOperations(words []string) bool {
	matchValidOperationPattern := regexp.MustCompile(`plus\s|minus\s|divided\sby\s|multiplied\sby\s`)

	for _, word := range words {
		if !matchValidOperationPattern.MatchString(word) {
			return false
		}
	}

	return true
}

func isNonMathQuestionOrUnsupportedOperation(expr string) bool {
	if err := isNonMathQuestion(expr); err != nil {
		return true
	}

	if err := isUnsupportedOperations(expr); err != nil {
		return true
	}

	return false
}

func NewUnsupportedExpressionError(expr string) ExpressionError {
	return ExpressionError{
		Msg:        "Expression contains unsupported operations",
		Expression: expr,
		ErrType:    UnsupportedOperations,
	}
}

func NewInvalidSyntaxExpressionError(expr, msg string) ExpressionError {
	return ExpressionError{
		Msg:        msg,
		Expression: expr,
		ErrType:    InvalidSyntax,
	}
}
