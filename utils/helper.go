package utils

import "github.com/go-playground/validator/v10"

func Validate(r interface{}) validator.ValidationErrors {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func Pick[T any](newVal *T, old T) T {
	if newVal != nil {
		return *newVal
	}
	return old
}
