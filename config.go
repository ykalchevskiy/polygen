package main

import (
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const defaultDiscriminator = "type"

const (
	JSONVersionV1           = "v1"
	JSONVersionV2Experiment = "v2experiment"
	JSONVersionV2           = "v2"
)

// Config and TypeMapping hold normalized internal generator state.
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
	TypeNameMethod     string
}

type TypeMapping struct {
	SubType   string
	TypeName  string
	IsPointer bool
}

// FileConfig, FileTypeConfig, and FileSubtypeConfig represent the external JSON configuration.
type FileConfig struct {
	Types                   []FileTypeConfig `json:"types"`
	StrictByDefault         bool             `json:"strictByDefault,omitempty"`
	PointerByDefault        bool             `json:"pointerByDefault,omitempty"`
	DefaultDiscriminator    string           `json:"defaultDiscriminator,omitempty"`
	DefaultBuildTag         string           `json:"defaultBuildTag,omitempty"`
	DefaultJSONVersion      string           `json:"defaultJsonVersion,omitempty"`
	TypeNameMethodByDefault string           `json:"typeNameMethodByDefault,omitempty"`
}

type FileTypeConfig struct {
	Type           string                       `json:"type"`
	Interface      string                       `json:"interface"`
	Package        string                       `json:"package"`
	Subtypes       map[string]FileSubtypeConfig `json:"subtypes"`
	Directory      string                       `json:"directory,omitempty"`
	Filename       string                       `json:"filename,omitempty"`
	Discriminator  string                       `json:"discriminator,omitempty"`
	Strict         *bool                        `json:"strict,omitempty"`
	DefaultSubtype string                       `json:"defaultSubtype,omitempty"`
	BuildTag       string                       `json:"buildTag,omitempty"`
	JSONVersion    string                       `json:"jsonVersion,omitempty"`
	TypeNameMethod *string                      `json:"typeNameMethod,omitempty"`
}

type FileSubtypeConfig struct {
	Name    *string `json:"name,omitempty"`
	Pointer *bool   `json:"pointer,omitempty"`
}

func convertFileConfigToConfig(typeConfig *FileTypeConfig, config *FileConfig) *Config {
	cfg := &Config{
		Type:           typeConfig.Type,
		Interface:      typeConfig.Interface,
		Package:        typeConfig.Package,
		Discriminator:  typeConfig.Discriminator,
		Strict:         config.StrictByDefault,
		BuildTag:       config.DefaultBuildTag,
		JSONVersion:    config.DefaultJSONVersion,
		TypeNameMethod: config.TypeNameMethodByDefault,
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

	if v := cfg.JSONVersion; v == "" {
		cfg.JSONVersion = JSONVersionV1
	} else if v != JSONVersionV1 && v != JSONVersionV2Experiment && v != JSONVersionV2 {
		cfg.JSONVersion = JSONVersionV1
	}

	if typeConfig.TypeNameMethod != nil {
		cfg.TypeNameMethod = *typeConfig.TypeNameMethod
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
