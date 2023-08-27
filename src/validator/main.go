package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	ValidateStruct(s any) error
}

type GoValidator struct {
	*validator.Validate
}

var _ Validator = (*GoValidator)(nil)

func (v *GoValidator) ValidateStruct(s any) error {
	return v.Struct(s)
}

func New() Validator {
	return &GoValidator{validator.New()}
}
