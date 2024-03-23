package util

import (
	"errors"
	"gocrud/internal/core/domain"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() (Validator, error) {
	validate := validator.New()

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return Validator{}, errors.New("failed to get translator")
	}

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return Validator{}, err
	}

	return Validator{validate: validate, trans: trans}, nil
}

func (v Validator) Validate(data interface{}) *domain.ErrorValidationResponse {
	if err := v.validate.Struct(data); err != nil {
		var errMessages []domain.ValidationError
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errMessages = append(errMessages, domain.ValidationError{
				Field: e.Field(),
				Error: e.Translate(v.trans),
			})
		}
		return &domain.ErrorValidationResponse{
			Message: "validation error",
			Errors:  errMessages,
		}
	}

	return nil
}
