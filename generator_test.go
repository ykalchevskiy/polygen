package main

import (
	"bytes"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		isPointerTrue := true
		config := FileConfig{
			Types: []FileTypeConfig{
				{
					Type:      "TestType",
					Interface: "TestInterface",
					Package:   "test",
					Subtypes: map[string]FileSubtypeConfig{
						"SubType1": {},
						"SubType2": {
							Pointer: &isPointerTrue,
						},
					},
				},
			},
		}

		cfg := convertFileConfigToConfig(&config.Types[0], &config)

		// Generate code
		code, err := generate(cfg)
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		if code == nil {
			t.Error("generated code is empty")
		}

		// Test required components
		required := []string{
			"package test",
			"type TestType struct {",
			"TestInterface",
			"func (v TestType) MarshalJSON() ([]byte, error)",
			"func (v *TestType) UnmarshalJSON(data []byte) error",
			`"type", typeName`,
			`case "sub-type-1":`,
			`case "sub-type-2":`,
			`reflect.TypeOf((*SubType1)(nil)).Elem():`,
			`reflect.TypeOf((*SubType2)(nil)):`,
		}

		for _, r := range required {
			if !bytes.Contains(code, []byte(r)) {
				t.Errorf("generated code missing required part: %q", r)
				t.Logf("Generated code:\n%s", string(code))
			}
		}
	})

	t.Run("custom values", func(t *testing.T) {
		subType1Name := "my-subtype-1"
		config := FileConfig{
			DefaultDescriptor: "kind",
			StrictByDefault:   true,
			PointerByDefault:  true,
			Types: []FileTypeConfig{
				{
					Type:      "TestType",
					Interface: "TestInterface",
					Package:   "test",
					Directory: "pkg",
					Subtypes: map[string]FileSubtypeConfig{
						"SubType1": {
							Name: &subType1Name,
						},
						"SubType2": {},
					},
				},
			},
		}

		cfg := convertFileConfigToConfig(&config.Types[0], &config)

		// Generate code
		code, err := generate(cfg)
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}

		if code == nil {
			t.Error("generated code is empty")
		}

		// Test required components
		required := []string{
			"package test",
			"type TestType struct {",
			"TestInterface",
			"func (v TestType) MarshalJSON() ([]byte, error)",
			"func (v *TestType) UnmarshalJSON(data []byte) error",
			`"kind", typeName`,
			`case "my-subtype-1":`,
			`case "sub-type-2":`,
			`reflect.TypeOf((*SubType1)(nil)):`,
			`reflect.TypeOf((*SubType2)(nil)):`,
			"decoder.DisallowUnknownFields()",
		}

		for _, r := range required {
			if !bytes.Contains(code, []byte(r)) {
				t.Errorf("generated code missing required part: %q", r)
				t.Logf("Generated code:\n%s", string(code))
			}
		}
	})
}
