// https://blog.depa.do/post/gin-validation-errors-handling
// https://github.com/go-playground/validator/issues/258#issuecomment-257266610
package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func FormatErrors(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}
