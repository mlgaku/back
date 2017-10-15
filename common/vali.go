package common

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vl "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type vali struct {
	Validate  *vl.Validate
	Translate ut.Translator
}

// 初始化
func (v *vali) init() {
	v.Validate = vl.New()
	v.Translate, _ = ut.New(en.New()).GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.Validate, v.Translate)
}

// 验证 var
func (v *vali) Var(field interface{}, tag string) string {
	err := v.Validate.Var(field, tag)
	if err == nil {
		return ""
	}
	for _, err := range err.(vl.ValidationErrors) {
		return err.Translate(v.Translate)
	}
	return ""
}

// 验证 struct
func (v *vali) Struct(s interface{}) string {
	if err := v.Validate.Struct(s); err != nil {
		for _, err := range err.(vl.ValidationErrors) {
			return err.Translate(v.Translate)
		}
	}
	return ""
}

// 获得 vali 实例
func NewVali() *vali {
	v := &vali{}
	v.init()
	return v
}
