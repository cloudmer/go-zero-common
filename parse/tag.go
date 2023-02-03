package parse

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
	"reflect"
	"regexp"
)

// 自定义翻译 注册tag
func (ins *validate) registerTag()  {
	// 注册tag 错误信息 字段名称
	ins.Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		if tag := field.Tag.Get("json"); tag != "" {
			return tag
		}
		return field.Name
	})

	// 注册一个国家码验证器
	_ = ins.Validate.RegisterValidation("cc", func(fl validator.FieldLevel) bool {
		ok, _ := regexp.MatchString(`^[1-9][0-9]{1,2}$`, fl.Field().String())
		return ok
	})
	// 注册一个国家码验证器 翻译
	ins.tagTranslate(ins.Validate, "cc", map[string]string{
		language.English.String(): "{0} is a invalid cc.",
		language.Chinese.String(): "{0} 不是一个可用的国家码",
	})

}

// tag 翻译
func (ins *validate) tagTranslate(validate *validator.Validate, tag string, messages map[string]string) {
	for lang, message := range messages {
		_ = validate.RegisterTranslation(tag, ins.translators[lang], func(ut ut.Translator) error {
			return ut.Add(tag, message, false)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})
	}
}