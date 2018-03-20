package rest

import (
	"net/http"
)

type IHttpError interface {
	StatusCode() int
}
type BadRequest string

func (e BadRequest) Error() string {
	return string(e)
}

func (e BadRequest) StatusCode() int {
	return http.StatusBadRequest
}

func WrapBadRequest(err error, message string) error {
	if err != nil {
		return BadRequest(message + ":" + err.Error())
	}
	return nil
}

func BadRequestNotFound(err error) error {
	if err != nil {
		return NotFound(err.Error())
	}
	return nil
}

func BadRequestValid(err error) error {
	if err != nil {
		return ValidError(err.Error())
	}
	return nil
}

func IsErrorRecord(err error) error {
	if err != nil && err.Error() == "not found" {
		return nil
	}
	return err
}

func BadRequestValidString(err string) error {
	if len(err) > 0 {
		return ValidError(err)
	}
	return nil
}

type Unauthorized string

func (e Unauthorized) Error() string {
	return string(e)
}

func (e Unauthorized) StatusCode() int {
	return http.StatusUnauthorized
}

type InternalServerError string

func (e InternalServerError) Error() string {
	return string(e)
}

func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

type NotFound string

func (e NotFound) Error() string {
	return string(e)
}

func (e NotFound) StatusCode() int {
	return http.StatusNotFound
}

type ValidError string

func (e ValidError) Error() string {
	return string(e)
}

func (e ValidError) StatusCode() int {
	return http.StatusForbidden
}

func AssertNil(errs ...error) {
	for _, item := range errs {
		if item != nil {
			panic(item)
		}
	}
}
func IsNotFound(err error) bool {
	if err.Error() == "not found" {
		return true
	}
	return false
}
