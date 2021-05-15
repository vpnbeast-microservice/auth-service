package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	translator locales.Translator
	trans      ut.Translator
	uni        *ut.UniversalTranslator
	v          *validator.Validate
)

func init() {
	translator = en.New()
	uni = ut.New(translator, translator)

	var found bool
	trans, found = uni.GetTranslator("en")
	if !found {
		panic("can not find translator")
	}

	v = validator.New()
	err := entranslations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		panic(err)
	}
}

func isValidRequest(context *gin.Context, request *authRequest) (bool, []string) {
	var errSlice []string
	var err error
	if err = context.ShouldBindJSON(request); err == nil {
		if err = v.Struct(request); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				errSlice = append(errSlice, e.Translate(trans))
			}
			return false, errSlice
		}
		return true, errSlice
	}

	errSlice = append(errSlice, "not a valid JSON request!")
	return false, errSlice
}

// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go
/*func translate() {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must be greater", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} must be smaller", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field())
		return t
	})}
}*/
