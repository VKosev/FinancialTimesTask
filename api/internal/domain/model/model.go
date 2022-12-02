package model

type ExpressionError struct {
	expression string
	endpoint   string
	frequency  int
	errType    string
}
