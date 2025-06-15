package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"text/template"
)

//go:embed template.go.tmpl
var codeTemplate string

func generate(cfg *Config) (string, error) {
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to format code: %v", err)
	}

	return string(formatted), nil
}
