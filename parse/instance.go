package parse

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sync"
)

var (
	instance *validate
	once     sync.Once
)

type validate struct {
	*validator.Validate
	translators map[string]ut.Translator
}

// Verify 验证
func (ins *validate) Verify(req *http.Request, v interface{}) error {
	return ins.ErrTranslate(req.Context(), ins.Validate.Struct(v))
}

// GetInstance 单例模式
func GetInstance() *validate {
	once.Do(func() {
		instance = &validate{
			Validate: validator.New(),
			translators: make(map[string]ut.Translator),
		}
		// 翻译
		instance.translate()
		// tag 注册 自定义注册验证规则
		instance.registerTag()
	})
	return instance
}
