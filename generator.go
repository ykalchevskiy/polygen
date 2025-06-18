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

func generate(cfg *Config) ([]byte, error) {
	tmpl, err := template.New("code").Parse(codeTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return nil, fmt.Errorf("executing template: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("formatting code: %v", err)
	}

	return formatted, nil
}
