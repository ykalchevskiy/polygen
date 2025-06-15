package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

func main() {
	configPath := flag.String("config", ".polygen.json", "Path to the configuration file")
	flag.Parse()

	// Read and parse config file
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config FileConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Set default descriptor if not specified
	if config.DefaultDescriptor == "" {
		config.DefaultDescriptor = "type"
	}

	// Generate code for each type
	configDir := filepath.Dir(*configPath)
	for _, typeConfig := range config.Types {
		// Convert config to internal format
		cfg := &Config{
			Type:       typeConfig.Type,
			Interface:  typeConfig.Interface,
			Descriptor: typeConfig.Descriptor,
			Package:    typeConfig.Package,
		}

		// Use default descriptor if not specified
		if cfg.Descriptor == "" {
			cfg.Descriptor = config.DefaultDescriptor
		}

		// Set strict mode
		if typeConfig.Strict != nil {
			cfg.Strict = *typeConfig.Strict
		} else {
			cfg.Strict = config.StrictByDefault
		}

		// Process subtypes
		for subType, subConfig := range typeConfig.Subtypes {
			var typeName string
			if subConfig.Name != nil {
				typeName = *subConfig.Name
			} else {
				// Convert subtype name to kebab-case by default
				typeName = toKebabCase(subType)
			}
			cfg.Types = append(cfg.Types, TypeMapping{
				SubType:   subType,
				TypeName:  typeName,
				IsPointer: subConfig.Pointer,
			})
		}

		// Sort types by SubType for consistent output
		sort.Slice(cfg.Types, func(i, j int) bool {
			return cfg.Types[i].SubType < cfg.Types[j].SubType
		})

		// Set output path
		var outputPath string
		if typeConfig.Directory != "" {
			outputPath = filepath.Join(configDir, typeConfig.Directory)
		} else {
			outputPath = configDir
		}

		if typeConfig.Filename != "" {
			outputPath = filepath.Join(outputPath, typeConfig.Filename)
		} else {
			// Default filename uses snake_case to preserve consistency with _polygen.go suffix
			outputPath = filepath.Join(outputPath, toSnakeCase(typeConfig.Type)+"_polygen.go")
		}

		// Create output directory if it doesn't exist
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		// Generate code
		code, err := generate(cfg)
		if err != nil {
			log.Fatalf("Failed to generate code for type %s: %v", typeConfig.Type, err)
		}

		// Write generated code to file
		if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
			log.Fatalf("Failed to write generated code: %v", err)
		}
	}
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
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, strings.ToLower(s[lastPos:i]))
			lastPos = i
		}
	}

	// append the last word
	if lastPos < len(s) {
		words = append(words, strings.ToLower(s[lastPos:]))
	}

	result = strings.Join(words, sep)
	return result
}
