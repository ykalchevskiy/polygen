package main

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const defaultDiscriminator = "type"

// Config represents the internal configuration used by the generator
type Config struct {
	Type           string
	Interface      string
	Package        string
	Types          []TypeMapping
	Discriminator  string
	Strict         bool
	DefaultSubtype string
}

// TypeMapping represents a mapping between a concrete type and its JSON type name
type TypeMapping struct {
	SubType   string
	TypeName  string
	IsPointer bool
}

// FileConfig represents the configuration file structure
type FileConfig struct {
	// Types is a list of type configurations to generate
	Types []FileTypeConfig `json:"types"`
	// StrictByDefault determines if strict mode should be enabled by default for all types
	StrictByDefault bool `json:"strictByDefault,omitempty"`
	// PointerByDefault determines if pointer mode should be enabled by default for all subtypes
	PointerByDefault bool `json:"pointerByDefault,omitempty"`
	// DefaultDiscriminator is the default JSON field name to distinguish types
	DefaultDiscriminator string `json:"defaultDiscriminator,omitempty"`
}

// FileTypeConfig represents configuration for a single polymorphic type
type FileTypeConfig struct {
	// Type is the name of the polymorphic structure to generate
	Type string `json:"type"`
	// Interface is the name of the interface all subtypes implement
	Interface string `json:"interface"`
	// Package is the package name for the generated file
	Package string `json:"package"`
	// Subtypes maps Go types to their configurations
	Subtypes map[string]FileSubtypeConfig `json:"subtypes"`
	// Directory is the output directory path, relative to the config file
	Directory string `json:"directory,omitempty"`
	// Filename is the output filename, defaults to <type>_polygen.go in snake_case
	Filename string `json:"filename,omitempty"`
	// Discriminator is the JSON field name to distinguish types
	Discriminator string `json:"discriminator,omitempty"`
	// Strict enables strict JSON unmarshaling for this type
	Strict *bool `json:"strict,omitempty"`
	// DefaultSubtype is the default subtype to use when the discriminator is missing
	DefaultSubtype string `json:"defaultSubtype,omitempty"`
}

// FileSubtypeConfig represents configuration for a subtype
type FileSubtypeConfig struct {
	// Name is the JSON type name, defaults to the subtype name in snake_case if not specified
	Name *string `json:"name,omitempty"`
	// Pointer indicates if this type should be used as a pointer
	Pointer *bool `json:"pointer,omitempty"`
}

func convertFileConfigToConfig(typeConfig *FileTypeConfig, config *FileConfig) *Config {
	cfg := &Config{
		Type:           typeConfig.Type,
		Interface:      typeConfig.Interface,
		Package:        typeConfig.Package,
		Discriminator:  typeConfig.Discriminator,
		Strict:         config.StrictByDefault,
		DefaultSubtype: typeConfig.DefaultSubtype,
	}

	if cfg.Discriminator == "" {
		cfg.Discriminator = config.DefaultDiscriminator
	}
	if cfg.Discriminator == "" {
		cfg.Discriminator = defaultDiscriminator
	}

	if typeConfig.Strict != nil {
		cfg.Strict = *typeConfig.Strict
	}

	for subType, subCfg := range typeConfig.Subtypes {
		var typeName string
		if subCfg.Name != nil {
			typeName = *subCfg.Name
		} else {
			typeName = toKebabCase(subType)
		}

		var isPointer bool
		if subCfg.Pointer != nil {
			isPointer = *subCfg.Pointer
		} else {
			isPointer = config.PointerByDefault
		}

		cfg.Types = append(cfg.Types, TypeMapping{
			SubType:   subType,
			TypeName:  typeName,
			IsPointer: isPointer,
		})
	}

	// Sort types by SubType for consistent output
	sort.Slice(cfg.Types, func(i, j int) bool {
		return cfg.Types[i].SubType < cfg.Types[j].SubType
	})

	return cfg
}

func getOutputPath(typeConfig *FileTypeConfig, configDir string) string {
	var outputPath string
	if typeConfig.Directory != "" {
		outputPath = filepath.Join(configDir, typeConfig.Directory)
	} else {
		outputPath = configDir
	}

	if typeConfig.Filename != "" {
		outputPath = filepath.Join(outputPath, typeConfig.Filename)
	} else {
		outputPath = filepath.Join(outputPath, toSnakeCase(typeConfig.Type)+"_polygen.go")
	}

	return outputPath
}

// toKebabCase converts a string from PascalCase to kebab-case
func toKebabCase(s string) string {
	return toCase(s, rune('-'))
}

// toSnakeCase converts a string from PascalCase to snake_case
func toSnakeCase(s string) string {
	return toCase(s, rune('_'))
}

func toCase(s string, sep rune) string {
	var result strings.Builder

	runes := []rune(s)
	length := len(runes)

	for i := 0; i < length; i++ {
		current := runes[i]

		if i > 0 {
			prev := runes[i-1]
			next := rune(0)
			if i+1 < length {
				next = runes[i+1]
			}

			// lower -> Upper (e.g., myTest -> my_test)
			if unicode.IsLower(prev) && unicode.IsUpper(current) {
				result.WriteRune(sep)
			} else
			// (e.g., HTTPServer -> http_server)
			if unicode.IsUpper(prev) && unicode.IsUpper(current) && next != 0 && unicode.IsLower(next) {
				result.WriteRune(sep)
			} else
			// (e.g., test123 -> test_123)
			if (unicode.IsLetter(prev) && unicode.IsDigit(current)) ||
				(unicode.IsDigit(prev) && unicode.IsLetter(current)) {
				result.WriteRune(sep)
			}
		}

		result.WriteRune(unicode.ToLower(current))
	}

	return result.String()
}
