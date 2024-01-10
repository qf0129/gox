package validx

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/qf0129/gox/errx"
)

var once sync.Once
var _instance *validator.Validate

func instance() *validator.Validate {
	once.Do(func() {
		_instance = validator.New()
	})
	return _instance
}

func Validate(data any) *errx.Err {
	err := instance().Struct(data)
	if err != nil {
		if errList, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errList {
				return errx.ValidateParamFailed.AddMsg(e.Field())
			}
		} else {
			return errx.ValidateParamFailed.AddErr(err)
		}
	}
	return nil
}
