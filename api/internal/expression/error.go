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

type ExpressionError struct {
	Msg        string
	Expression string
	ErrType    ErrorType
}

func (e ExpressionError) Error() string {
	return e.Msg
}

func IsExpressionError(err error) (ExpressionError, bool) {
	e, ok := err.(ExpressionError)
	if ok {
		return e, true
	}

	return ExpressionError{}, false
}

type ErrorHistory struct {
	Expression string
	Endpoints  map[string]int
	Frequency  int
	ErrType    ErrorType
}

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
