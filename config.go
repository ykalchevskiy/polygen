package main

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const defaultDiscriminator = "type"

const (
	JSONVersionV1   = "v1"
	JSONVersionV2   = "v2"
	JSONVersionBoth = "both"
)

// Config represents the internal configuration used by the generator.
type Config struct {
	Type               string
	Interface          string
	Package            string
	Types              []TypeMapping
	Discriminator      string
	Strict             bool
	DefaultSubtypeName string
	BuildTag           string
	JSONVersion        string
}

// TypeMapping represents a mapping between a concrete type and its JSON type name.
type TypeMapping struct {
	SubType   string
	TypeName  string
	IsPointer bool
}

// FileConfig represents the configuration file structure.
type FileConfig struct {
	// Types is a list of type configurations to generate
	Types []FileTypeConfig `json:"types"`
	// StrictByDefault determines if strict mode should be enabled by default for all types (does not apply to jsonv2)
	StrictByDefault bool `json:"strictByDefault,omitempty"`
	// PointerByDefault determines if pointer mode should be enabled by default for all subtypes
	PointerByDefault bool `json:"pointerByDefault,omitempty"`
	// DefaultDiscriminator is the default JSON field name to distinguish types
	DefaultDiscriminator string `json:"defaultDiscriminator,omitempty"`
	// DefaultBuildTag is the build constraint for generated code
	DefaultBuildTag string `json:"defaultBuildTag,omitempty"`
	// JSONVersionByDefault determines the json version generation enabled by default (v1, v2, both)
	JSONVersionByDefault string `json:"jsonVersionByDefault,omitempty"`
}

// FileTypeConfig represents configuration for a single polymorphic type.
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
	// Strict enables strict JSON unmarshaling for this type (does not apply to jsonv2)
	Strict *bool `json:"strict,omitempty"`
	// DefaultSubtype is the default subtype to use when the discriminator is missing
	DefaultSubtype string `json:"defaultSubtype,omitempty"`
	// BuildTag is the build constraint for this type
	BuildTag string `json:"buildTag,omitempty"`
	// JSONVersion enables generation of jsonv2 code for this type (v1, v2, both)
	JSONVersion string `json:"jsonVersion,omitempty"`
}

// FileSubtypeConfig represents configuration for a subtype.
type FileSubtypeConfig struct {
	// Name is the JSON type name, defaults to the subtype name in snake_case if not specified
	Name *string `json:"name,omitempty"`
	// Pointer indicates if this type should be used as a pointer
	Pointer *bool `json:"pointer,omitempty"`
}

func convertFileConfigToConfig(typeConfig *FileTypeConfig, config *FileConfig) *Config {
	cfg := &Config{
		Type:          typeConfig.Type,
		Interface:     typeConfig.Interface,
		Package:       typeConfig.Package,
		Discriminator: typeConfig.Discriminator,
		Strict:        config.StrictByDefault,
		BuildTag:      config.DefaultBuildTag,
		JSONVersion:   config.JSONVersionByDefault,
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

	if typeConfig.BuildTag != "" {
		cfg.BuildTag = typeConfig.BuildTag
	}

	if typeConfig.JSONVersion != "" {
		cfg.JSONVersion = typeConfig.JSONVersion
	}

	if cfg.JSONVersion == "" {
		cfg.JSONVersion = JSONVersionV1
	} else if cfg.JSONVersion != JSONVersionV1 && cfg.JSONVersion != JSONVersionV2 && cfg.JSONVersion != JSONVersionBoth {
		cfg.JSONVersion = JSONVersionV1
	}

	var defaultSubtypeName string

	for subType, subCfg := range typeConfig.Subtypes {
		var typeName string
		if subCfg.Name != nil {
			typeName = *subCfg.Name
		} else {
			typeName = toKebabCase(subType)
		}

		if subType == typeConfig.DefaultSubtype {
			defaultSubtypeName = typeName
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

	cfg.DefaultSubtypeName = defaultSubtypeName

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

// toKebabCase converts a string from PascalCase to kebab-case.
func toKebabCase(s string) string {
	return toCase(s, rune('-'))
}

// toSnakeCase converts a string from PascalCase to snake_case.
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

			switch {
			case unicode.IsLower(prev) && unicode.IsUpper(current):
				// lower -> Upper (e.g., myTest -> my_test)
				result.WriteRune(sep)
			case unicode.IsUpper(prev) && unicode.IsUpper(current) && next != 0 && unicode.IsLower(next):
				// (e.g., HTTPServer -> http_server)
				result.WriteRune(sep)
			case unicode.IsLetter(prev) && unicode.IsDigit(current):
				// (e.g., test123 -> test_123)
				result.WriteRune(sep)
			case unicode.IsDigit(prev) && unicode.IsLetter(current):
				// (e.g., version2Test -> version_2_test)
				result.WriteRune(sep)
			}
		}

		result.WriteRune(unicode.ToLower(current))
	}

	return result.String()
}
