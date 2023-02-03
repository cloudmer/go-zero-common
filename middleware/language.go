package middleware

import (
	"context"
	"net/http"
)

// LanguageMiddleware 中间件 获取请求翻译语言
func LanguageMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// 从 get参数 获取
		languageStr := request.URL.Query().Get("language")
		// 从 header头 获取
		if languageStr == "" {
			languageStr = request.Header.Get("Accept-Language")
		}
		// 设置上下文 请求翻译的语言
		ctx := context.WithValue(request.Context(), "language", languageStr)
		next.ServeHTTP(writer, request.WithContext(ctx))
	}
}