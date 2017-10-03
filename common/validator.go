package common

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vl "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type validator struct {
	Validate  *vl.Validate
	Translate ut.Translator
}

// 初始化
func (v *validator) init() {
	v.Validate = vl.New()
	v.Translate, _ = ut.New(en.New()).GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.Validate, v.Translate)
}

// 验证 struct
func (v *validator) Struct(s interface{}) string {
	if err := v.Validate.Struct(s); err != nil {
		for _, err := range err.(vl.ValidationErrors) {
			return err.Translate(v.Translate)
		}
	}
	return ""
}

// 获得 validator 实例
func NewValidator() *validator {
	v := &validator{}
	v.init()
	return v
}
