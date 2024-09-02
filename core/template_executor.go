package core

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/go-logr/logr"
)

type TemplateExecutor struct {
	logger logr.Logger
}

func NewTemplateExecutor(logger logr.Logger) *TemplateExecutor {
	return &TemplateExecutor{logger: logger}
}

func (te *TemplateExecutor) Execute(variables []string) (string, error) {
	var result bytes.Buffer

	templates, err := templatesFS.ReadDir("templates")
	if err != nil {
		return "", fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, templateFile := range templates {
		templateContent, err := templatesFS.ReadFile(
			fmt.Sprintf("templates/%s", templateFile.Name()),
		)
		if err != nil {
			return "", fmt.Errorf("failed to read template file %s: %w", templateFile.Name(), err)
		}

		tmpl, err := template.New(templateFile.Name()).Parse(string(templateContent))
		if err != nil {
			return "", fmt.Errorf("failed to parse template %s: %w", templateFile.Name(), err)
		}

		result.WriteString(fmt.Sprintf("# %s\n", templateFile.Name()))

		err = tmpl.Execute(&result, map[string]interface{}{"Variables": variables})
		if err != nil {
			return "", fmt.Errorf("failed to execute template %s: %w", templateFile.Name(), err)
		}

		result.WriteString("\n")
	}

	return result.String(), nil
}

func (te *TemplateExecutor) ExecuteToWriter(variables []string, writer io.Writer) error {
	result, err := te.Execute(variables)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(result))
	return err
}
