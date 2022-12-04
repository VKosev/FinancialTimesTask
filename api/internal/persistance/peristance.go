package persistance

import (
	"errors"
	"strings"
	"unicode"

	"github.com/vkosev/ft_api/internal/expression"
)

// Connection is persistance connection used to store and retrieve expression errors
// and implements expression.ErrorPersister interface.
type Connection struct {
	errors map[string]errorEntity
}

// errorEntity is the type the connection package uses for persisting and retrieving expression errors.
type errorEntity struct {
	Expression string
	Endpoints  map[string]int
	Frequency  int
	ErrType    expression.ErrorType
}

// NewConnection creates a persistance connection.
func NewConnection() *Connection {
	return &Connection{
		errors: make(map[string]errorEntity),
	}
}

// Save persists the passed expression.ExpressionError and url.
// Returns an error if problem occurs.
func (c *Connection) Save(err expression.ExpressionError, url string) error {
	if err := fieldsAreNotNull(err, url); err != nil {
		return err
	}

	expr := trimWhiteSpaces(err.Expression)

	existingErr, ok := c.errors[expr]
	if ok {
		existingErr.Frequency++
		addEndpoint(existingErr.Endpoints, url)
		c.errors[expr] = existingErr

		return nil
	}

	newError := errorEntity{
		Expression: err.Expression,
		Frequency:  1,
		ErrType:    err.ErrType,
		Endpoints: map[string]int{
			url: 1,
		},
	}
	c.errors[expr] = newError

	return nil
}

func (c *Connection) GetAll() []expression.ErrorHistory {
	var errors []expression.ErrorHistory

	for _, err := range c.errors {
		historyErr := expression.ErrorHistory{
			Expression: err.Expression,
			Frequency:  err.Frequency,
			ErrType:    err.ErrType,
			Endpoints:  err.Endpoints,
		}

		errors = append(errors, historyErr)
	}

	return errors
}

// addEndpoint inserts the url into the map if key already exist
// otherwise it creates new element in the map.
func addEndpoint(endpoints map[string]int, url string) {
	count, ok := endpoints[url]

	if ok {
		count++
		endpoints[url] = count
	} else {
		endpoints = map[string]int{
			url: 1,
		}
	}
}

// fieldsAreNotNull validates that the fields of the passed expression.ExpressionError
// and url are not null
func fieldsAreNotNull(err expression.ExpressionError, url string) error {
	if err.Expression == "" {
		return errors.New("error expression is empty")
	}

	if url == "" {
		return errors.New("endpoint url is empty")
	}

	return nil
}

// trimWhiteSpace removes all  white spaces of the passed string
// and returns new string containing the passed string without white spaces
func trimWhiteSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
