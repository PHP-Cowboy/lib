package rds

import "github.com/go-playground/validator/v10"

func InitValidate() *validator.Validate {
	return validator.New()
}
