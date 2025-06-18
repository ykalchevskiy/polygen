package main

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const defaultDescriptor = "type"

// Config represents the internal configuration used by the generator
type Config struct {
	Type       string
	Interface  string
	Package    string
	Types      []TypeMapping
	Descriptor string
	Strict     bool
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
	// DefaultDescriptor is the default JSON field name to distinguish types
	DefaultDescriptor string `json:"defaultDescriptor,omitempty"`
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
	// Descriptor is the JSON field name to distinguish types
	Descriptor string `json:"descriptor,omitempty"`
	// Strict enables strict JSON unmarshaling for this type
	Strict *bool `json:"strict,omitempty"`
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
		Type:       typeConfig.Type,
		Interface:  typeConfig.Interface,
		Package:    typeConfig.Package,
		Descriptor: typeConfig.Descriptor,
		Strict:     config.StrictByDefault,
	}

	if cfg.Descriptor == "" {
		cfg.Descriptor = config.DefaultDescriptor
	}
	if cfg.Descriptor == "" {
		cfg.Descriptor = defaultDescriptor
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
	return toCase(s, "-")
}

// toSnakeCase converts a string from PascalCase to snake_case
func toSnakeCase(s string) string {
	return toCase(s, "_")
}

func toCase(s, sep string) string {
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && !unicode.IsLower(rs[i]) {
			words = append(words, strings.ToLower(string(rs[lastPos:i])))
			lastPos = i
		}
	}

	// append the last word
	if lastPos < len(rs) {
		words = append(words, strings.ToLower(string(rs[lastPos:])))
	}

	result = strings.Join(words, sep)
	return result
}
