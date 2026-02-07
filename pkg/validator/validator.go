package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器
var Trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改 gin 框架中的 Validator 引擎属性，实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取 json tag 的自定义方法，让错误信息中的字段名显示为 json 定义的名称
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用语言，后面参数是支持的语言
		uni := ut.New(enT, zhT, enT)

		// 根据参数获取对应的翻译器
		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册对应语言的翻译业务逻辑
		switch locale {
		case "en":
			err = en_translations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(v, Trans)
		default:
			err = en_translations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}
