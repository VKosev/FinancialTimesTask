package expression

import (
	"regexp"
	"strings"
)

const msgNotValidExpression string = "Expression is not valid"

var startAsExpressionPattern *regexp.Regexp = regexp.MustCompile(`What is \d+ `)
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
	startWithMathQuestionPattern := regexp.MustCompile(`What is \d+ (plus|multiplied|minus|divided) \d+`)
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
	containsSupportedOperationsPattern := regexp.MustCompile(`plus \d+|minus \d+|multiplied by \d+|divided by \d+`)

	if startAsExpressionPattern.MatchString(expr) && endsWithNumberAndQuestionMarkPattern.MatchString(expr) {
		if !containsSupportedOperationsPattern.MatchString(expr) {
			return NewUnsupportedExpressionError(expr)
		}
	}

	finalOperationIsSupported := regexp.MustCompile(`multiplied by \d+\?$|divided by \d+\?$|plus \d+\?$|minus \d+\?$`)
	if containsSupportedOperationsPattern.MatchString(expr) && !finalOperationIsSupported.MatchString(expr) {
		return NewUnsupportedExpressionError(expr)
	}

	return nil
}

func isInvalidSyntaxError(expr string) error {
	matchValidExpressionPattern := regexp.MustCompile(`^What is \d+| plus \d+| minus \d+| divided by \d+| multiplied by \d+|\?$`)

	stringMatches := matchValidExpressionPattern.FindAllString(expr, -1)

	compareExpr := strings.Join(stringMatches, "")
	if expr != compareExpr {
		return NewInvalidSyntaxExpressionError(expr, msgNotValidExpression)
	}

	return nil
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
