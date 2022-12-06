package expression

import "regexp"

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
	if !isNonMathQuestionOrUnsupportedOperation(expr) {
		validAndInvalidOperations := regexp.MustCompile(`plus[a-zA-Z a-zA-Z]*|minus[a-zA-Z a-zA-Z]*|multiplied by[a-zA-Z a-zA-Z]*|divided by[a-zA-Z a-zA-Z]*`)

		matchedStrings := validAndInvalidOperations.FindAllString(expr, -1)

		if len(matchedStrings) == 0 {
			return NewInvalidSyntaxExpressionError(expr, "Expression does not contain operations")
		}

		if !eachOperationEndsWithNumber(expr) {
			return NewInvalidSyntaxExpressionError(expr, msgNotValidExpression)
		}

		if !containsOnlySupportedOperations(matchedStrings) {
			return NewInvalidSyntaxExpressionError(expr, msgNotValidExpression)
		}
	}

	return nil
}

func containsOnlySupportedOperations(words []string) bool {
	matchValidOperationPattern := regexp.MustCompile(`plus |minus |divided by |multiplied by `)

	for _, word := range words {
		if !matchValidOperationPattern.MatchString(word) {
			return false
		}
	}

	for i := 0; i < len(words); i++ {
		if matchValidOperationPattern.MatchString(words[i]) && matchValidOperationPattern.MatchString(words[i+1]) {

			return false
		}
	}

	return true
}

func eachOperationEndsWithNumber(expr string) bool {
	eachOperationEndsWithNumberPattern := regexp.MustCompile(`plus\s\d+|minus\s\d+|divided\sby\s\d+|multiplied\sby\s\d+`)

	operations := eachOperationEndsWithNumberPattern.FindAllString(expr, -1)

	stringEndsWithNumberPattern := regexp.MustCompile(`\d+$`)

	for _, operation := range operations {
		if !stringEndsWithNumberPattern.MatchString(operation) {
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
