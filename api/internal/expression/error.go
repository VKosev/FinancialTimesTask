package expression

import (
	"encoding/json"
	"errors"
)

type ErrorType int

const (
	UnsupportedOperations ErrorType = iota
	NonMathQuestion
	InvalidSyntax
)

// ErrorPersister is used to persist errors from package expression
type ErrorPersister interface {
	// Save persists an expression.ExpressionError and the url,
	// returns an error if persisting fails.
	Save(exprError ExpressionError, url string) error

	// GetAll returns all persisted errors as slice of expression.ErrorHistory.
	GetAll() []ErrorHistory
}

// ExpressionError is returned for expression related problems.
type ExpressionError struct {
	Msg        string
	Expression string
	ErrType    ErrorType
}

func (e ExpressionError) Error() string {
	return e.Msg
}

// IsExpressionError validates if error can be casted to ExpressionError.
//
// Returns the successfully casted ExpressionError and true if sucessfull
// otherwise nil and false.
func IsExpressionError(err error) (ExpressionError, bool) {
	e, ok := err.(ExpressionError)
	if ok {
		return e, true
	}

	return ExpressionError{}, false
}

// ErrorHistory is returned when retrieving persisted errors.
type ErrorHistory struct {
	// The expression for which the error occured.
	Expression string
	// Endpoints map represents the urls on which the error occured.
	//
	// key - represents the url
	// value - represents the number of times this url with this expression recieved an error
	Endpoints map[string]int
	// Number of times an error occured for the given expression
	Frequency int
	// The error type
	ErrType ErrorType
}

// Implement the stringer inferface. Defaults to "Invalid Syntax"
func (et ErrorType) String() string {
	switch et {
	case 0:
		return "Unsupported Operation"
	case 1:
		return "Non Math Question"
	case 2:
		return "Invalid Syntax"
	default:
		return "Invalid Syntax"
	}
}

// ParseErrorType tries to parse a string into ErrorType.
//
// Returns the ErrorType and nil if successfull,
// otherwise the default ErrorType and error.
func ParseErrorType(s string) (ErrorType, error) {
	switch s {
	case "unsupportedOperations":
		return UnsupportedOperations, nil
	case "nonMathQuestion":
		return NonMathQuestion, nil
	case "invalidSyntax":
		return InvalidSyntax, nil
	default:
		return InvalidSyntax, errors.New("string is not value of ErrorType")
	}
}

func (et ErrorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (et *ErrorType) UnmarshalJSON(data []byte) (err error) {
	var errType string
	if err := json.Unmarshal(data, &errType); err != nil {
		return err
	}
	if *et, err = ParseErrorType(errType); err != nil {
		return err
	}
	return nil
}
