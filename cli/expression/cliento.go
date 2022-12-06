package expression

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

var ErrEvaluation = errors.New("expression can't be evaluated")

type Client struct {
	basePath string
}

// NewClient returns pointer to new instance of expression client.
//
// Accepts:
//
// - basePath: base url path, that will be prefixed to every request.
//	defaults to localhost:8080
func NewClient(basePath *string) *Client {
	return &Client{
		basePath: viper.GetString("clientApi.basePath"),
	}
}

func (c *Client) Evaluate(expr string) (*EvaluationResult, error) {
	jsonBody, err := json.Marshal(ExpressionRequest{Expression: expr})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/expression", c.basePath)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusInternalServerError {
		return nil, responseMsgAsError(res)
	}

	if res.StatusCode == http.StatusBadRequest {
		var resError *ExpressionErrorResponse

		if err = json.NewDecoder(res.Body).Decode(&resError); err != nil {
			return nil, err
		}

		return &EvaluationResult{
			Result:  0,
			Type:    resError.Type,
			Message: resError.Message,
		}, ErrEvaluation
	}

	var respBody EvaluatedExpressionResponse

	if err = json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	result := &EvaluationResult{
		Result:  respBody.Result,
		Type:    "",
		Message: "",
	}
	return result, nil
}

func (c *Client) Validate(expr string) (*ValidationResult, error) {
	jsonBody, err := json.Marshal(ExpressionRequest{Expression: expr})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/validate", c.basePath)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusInternalServerError {
		return nil, responseMsgAsError(res)
	}

	var respBody ValidationResult

	if err = json.NewDecoder(res.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return &respBody, nil

}

func (c *Client) Errors() ([]ErrorHistoryResponse, error) {
	url := fmt.Sprintf("%s/errors", c.basePath)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var errors *[]ErrorHistoryResponse

	if res.StatusCode == http.StatusBadRequest {
		return nil, responseMsgAsError(res)
	}

	if err = json.NewDecoder(res.Body).Decode(&errors); err != nil {
		return nil, err
	}

	return *errors, nil
}

func responseMsgAsError(res *http.Response) error {
	var resText *string

	err := json.NewDecoder(res.Body).Decode(resText)
	if err != nil {
		return err
	}

	return fmt.Errorf(err.Error())
}
