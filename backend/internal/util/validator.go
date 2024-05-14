package util

import (
	"errors"
	"regexp"

	"github.com/shironxn/blanknotes/internal/core/domain"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidator() (*Validator, error) {
	validate := validator.New()

	if err := validate.RegisterValidation("image", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`\.(?:jpg|png)$`).MatchString(fl.Field().String())
	}); err != nil {
		return nil, err
	}

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, errors.New("failed to get translator")
	}

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	return &Validator{validate: validate, trans: trans}, nil
}

func (v Validator) Validate(data interface{}) *domain.ErrorResponse {
	if err := v.validate.Struct(data); err != nil {
		var errMessages []domain.ValidationError
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			errMessages = append(errMessages, domain.ValidationError{
				Field: e.Field(),
				Error: e.Translate(v.trans),
			})
		}
		return &domain.ErrorResponse{
			Code:  400,
			Error: errMessages,
		}
	}

	return nil
}
