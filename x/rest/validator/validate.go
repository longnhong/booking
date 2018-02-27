package validator

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func Validate(s interface{}) error {
	validate = validator.New()
	return validate.Struct(s)
}

func ErrorChain(funcs ...(func(params ...interface{}) error)) error {
	for _, item := range funcs {
		if item() != nil {
			return item()
		}
	}
	return nil
}
