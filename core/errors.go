package core

import (
	"fmt"
)

type TemplateError struct {
	TemplateName string
	Err          error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template error in %s: %v", e.TemplateName, e.Err)
}

func (e *TemplateError) Unwrap() error {
	return e.Err
}
