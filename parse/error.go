package parse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

// ErrTranslate 错误翻译
func (ins *validate) ErrTranslate(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// json 转换失败 格式不正确
	// SyntaxError是对JSON语法错误的描述。
	if _, ok := err.(*json.SyntaxError); ok {
		return errors.New("JSON语法错误")
	}


	// UnmarshalTypeError 描述的JSON值 不适用于特定Go类型的值
	// json解析失败，请检查格式是否正确，检查字符类型是否正确
	if obj, ok := err.(*json.UnmarshalTypeError); ok {
		return errors.New(fmt.Sprintf("%v 是一个%v类型", strings.ToLower(obj.Field), obj.Type))
	}
	// 解析成结构体 字符类型 错误 一般情况下都是浏览器录入值范围或类型造成的 断言
	num, ok := err.(*strconv.NumError)
	if ok {
		return errors.New(fmt.Sprintf("参数解析失败，解析方法：%s，解析值：%s，请检查值类型与范围！", num.Func, num.Num))
	}

	// 描述传递给“Struct”、“StructExcept”、“StructPartial”或“Field”的无效参数 断言
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return errors.New(fmt.Sprintf("param error: %s", invalid.Error()))
	}

	// 字段验证的详细错误信息 断言是否是 字段详细错误slice
	validationErrors, ok := err.(validator.ValidationErrors)

	// 默认翻译语言
	languageTag := language.Chinese
	languageStr := ctx.Value("language")
	if languageStr != nil {
		if langStr, ok := languageStr.(string); ok {
			if langTag := language.Make(langStr); langTag != language.Und {
				languageTag = langTag
			}
		}
	}

	for _, err := range validationErrors {
		return errors.New(err.Translate(GetInstance().translators[languageTag.String()]))
	}
	return nil
}
