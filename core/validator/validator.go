package validator

import (
	"team-flow/core/errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) *errors.AppError {
	if err := validate.Struct(s); err != nil {
		return errors.NewAppError(
			errors.ErrCodeValidation,
			"Validation failed",
			err,
		)
	}

	return nil
}
