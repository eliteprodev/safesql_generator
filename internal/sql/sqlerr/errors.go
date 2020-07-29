package sqlerr

import (
	"errors"
	"fmt"
)

var Exists = errors.New("already exists")
var NotFound = errors.New("does not exist")
var NotUnique = errors.New("is not unique")

type Error struct {
	Err      error
	Code     string
	Message  string
	Location int
	Line     int
	Column   int
	// Hint     string
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s %s", e.Message, e.Err.Error())
	} else {
		return e.Message
	}
}

func ColumnExists(rel, col string) *Error {
	return &Error{
		Err:     Exists,
		Code:    "42701",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\"", col, rel),
	}
}

func ColumnNotFound(rel, col string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42703",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\"", col, rel),
	}
}

func RelationExists(rel string) *Error {
	return &Error{
		Err:     Exists,
		Code:    "42P07",
		Message: fmt.Sprintf("relation \"%s\"", rel),
	}
}

func RelationNotFound(rel string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42P01",
		Message: fmt.Sprintf("relation \"%s\"", rel),
	}
}

func SchemaExists(name string) *Error {
	return &Error{
		Err:     Exists,
		Code:    "42P06",
		Message: fmt.Sprintf("schema \"%s\"", name),
	}
}

func SchemaNotFound(sch string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "3F000",
		Message: fmt.Sprintf("schema \"%s\"", sch),
	}
}

func TypeExists(typ string) *Error {
	return &Error{
		Err:     Exists,
		Code:    "42710",
		Message: fmt.Sprintf("type \"%s\"", typ),
	}
}

func TypeNotFound(typ string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42704",
		Message: fmt.Sprintf("type \"%s\"", typ),
	}
}

func FunctionNotFound(fun string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42704",
		Message: fmt.Sprintf("function \"%s\"", fun),
	}
}

func FunctionNotUnique(fn string) *Error {
	return &Error{
		Err:     NotUnique,
		Message: fmt.Sprintf("function name \"%s\"", fn),
	}
}
