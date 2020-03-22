package services

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/go-playground/locales/en"
	
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidatorService struct {
	model interface{}
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

//handler validator error
func (v *ValidatorService) HandlerError() error {
	err := validator.New().Struct(v.model)
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	
	validate = validator.New()
	_ = en_translations.RegisterDefaultTranslations(validate, trans)
	var str string
	if _, ok := err.(validator.ValidationErrors); ok {
		errs := err.(validator.ValidationErrors)
		_ = seelog.Error(errs)
		for k, v := range errs.Translate(trans) {
			str += k + "," + v
		}
	}
	if str == "" {
		return nil
	} else {
		return errors.New(str)
	}
}
