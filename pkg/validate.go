package pkg

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(ctx context.Context, obj any) error {
	var validate = validator.New()
	if err := validate.StructCtx(ctx, obj); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return validationErrors
		}
		return err
	}
	return nil
}
