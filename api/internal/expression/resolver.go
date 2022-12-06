package expression

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

const (
	matchNumbersRegexExpression = `\d+`
)

// Resolver has methods to evaluate arithmetic expressions from text
type Resolver struct {
	l  *log.Logger
	ep ErrorPersister
}

// NewResolver returns a pointer to new instance of Resolver.
func NewResolver(l *log.Logger, ep ErrorPersister) *Resolver {
	return &Resolver{
		l:  l,
		ep: ep,
	}
}

// Evaluate evaluates a text exression.
//
// returns:
//		int if evaluation is sucessfull
//		error which can be casted to ExpressionError,
//		for expression related error.
// 		Normal error for other issues.
func (r *Resolver) Evaluate(expr string, url string) (int, error) {
	if err := isValid(expr); err != nil {
		if e, ok := IsExpressionError(err); ok {
			r.ep.Save(e, url)
			return 0, e
		}
		r.l.Println("ERROR: ", err.Error())
		return 0, err
	}

	if isValidExpressionWithNoOperations(expr) {
		return r.findNumberNoOperations(expr), nil
	}

	value, err := r.resolveOperations(expr)
	if err != nil {
		r.l.Println("ERROR: ", err.Error())
		return 0, err
	}

	return value, nil
}

func (r *Resolver) Validate(expr string, url string) error {
	if err := isValid(expr); err != nil {

		if e, ok := IsExpressionError(err); ok {

			if err := r.ep.Save(e, url); err != nil {
				return err
			}

			return e
		}

		return err
	}

	return nil
}

func (r *Resolver) ErrorHistory() []ErrorHistory {
	return r.ep.GetAll()
}

// resolveOperations evaluates simple or complex expressions and returns
// the evaluated value as integer. Returns an error if problem occurs.
func (r *Resolver) resolveOperations(expr string) (int, error) {
	numbers, err := findAllNumbers(expr)
	if err != nil {
		r.l.Println("ERROR: ", err.Error())
		return 0, err
	}

	operations := findAllOperations(expr)
	if len(numbers) != (len(operations) + 1) {
		r.l.Println("ERROR: the numbers of operations or numbers incorrect")
		return 0, errors.New("the numbers of operations or numbers incorrect")
	}

	var result int
	for index, operation := range operations {
		if index == 0 {
			result, err = calculateArithmeticOperation(operation, numbers[index], numbers[index+1])
			if err != nil {
				return 0, err
			}
		} else {
			result, err = calculateArithmeticOperation(operation, result, numbers[index+1])
			if err != nil {
				return 0, err
			}
		}
	}

	return result, nil
}

// calculateArithmeticOperation tries to calculate v1 and v2 based on the passed operation.
//
// Returns the calculated integere and nil on success, otherwise 0 and error.
func calculateArithmeticOperation(operation string, v1, v2 int) (int, error) {
	switch operation {
	case "plus":
		return v1 + v2, nil
	case "minus":
		return v1 - v2, nil
	case "multiplied":
		return v1 * v2, nil
	case "divided":
		return v1 / v2, nil
	default:
		return 0, fmt.Errorf("unsupported arithmetic operation %s", operation)
	}
}

// returns all found allowed operation from the passed string.
func findAllOperations(expr string) []string {
	pattern := regexp.MustCompile("plus|multiplied|minus|divided")

	return pattern.FindAllString(expr, -1)
}

// findAllNumbers finds all numbers in the passed string
// and returns them as slice of ints
func findAllNumbers(expr string) ([]int, error) {
	pattern := regexp.MustCompile(`\d+`)

	stringNumbers := pattern.FindAllString(expr, -1)

	intNumbers, err := convertToIntegers(stringNumbers)
	if err != nil {
		return nil, err
	}

	return intNumbers, nil

}

// convertToIntegers converts slice of numbers which are string
// into slice of the same numbers converted to int,
// returns an error if something fails
func convertToIntegers(snumbers []string) ([]int, error) {
	convertedNumbers := make([]int, len(snumbers))

	for index, snum := range snumbers {
		convertedNumber, err := strconv.Atoi(snum)
		if err != nil {
			return nil, err
		}

		convertedNumbers[index] = convertedNumber
	}

	return convertedNumbers, nil
}

// findNumberNoOperations finds the first occurrence of number
// in the passed string and returns it as int
func (r *Resolver) findNumberNoOperations(expr string) int {
	findNumbersPattern := regexp.MustCompile(matchNumbersRegexExpression)

	value := findNumbersPattern.FindString(expr)

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		r.l.Println("ERROR: ", err.Error())
		return 1
	}

	return valueInt
}

// isValidExpressionWithNoOperations validates that the passed string is
// valid expression, which contains no arithmetic operations.
// returns true if it's valid, otherwise false
func isValidExpressionWithNoOperations(expr string) bool {
	onlyNumberPattern := regexp.MustCompile(`What\sis\s\d+\?`)

	return onlyNumberPattern.Match([]byte(expr))
}
