package util

import (
	"errors"
	"fmt"

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
	en := en.New()
	uni := ut.New(en, en)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, errors.New("failed to get translator")
	}
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &Validator{validate: validate, trans: trans}, nil
}

func (v *Validator) Validate(data interface{}) error {
	err := v.validate.Struct(data)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			field := e.Field()
			tag := e.Tag()
			translatedTag, _ := v.trans.T(tag, field)
			return fmt.Errorf("%s: %s", field, translatedTag)
		}
	}
	return nil
}
