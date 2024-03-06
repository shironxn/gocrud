package util

import (
	"fmt"
	"gocrud/internal/core/domain"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(data interface{}) []domain.ValidationError {
	err := v.validate.Struct(data)

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.validate, trans)

	if err != nil {
		errors := make([]domain.ValidationError, 0)

		for _, validationErr := range err.(validator.ValidationErrors) {
			translatedErr := fmt.Errorf(validationErr.Translate(trans))
			errors = append(errors, domain.ValidationError{
				Field:  validationErr.StructField(),
				Errors: translatedErr.Error(),
			})
		}

		return errors
	}

	return nil
}
