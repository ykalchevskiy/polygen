package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Type       string
	Interface  string
	Types      []TypeMapping
	Descriptor string
	Strict     bool
	Package    string
	File       string
}

type TypeMapping struct {
	SubType   string
	TypeName  string
	IsPointer bool
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	code, err := generate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(cfg.File, []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() (*Config, error) {
	typeFlag := flag.String("type", "", "name of the polymorphic structure (required)")
	interfaceFlag := flag.String("interface", "", "name of the interface all subtypes should implement (required)")
	typesFlag := flag.String("types", "", "comma-separated list of subtypes and their type names (required)")
	descriptorFlag := flag.String("descriptor", "type", "name of the JSON field to distinguish types")
	strictFlag := flag.Bool("strict", false, "enable strict JSON unmarshaling (disallow unknown fields)")
	packageFlag := flag.String("package", "", "package name (defaults to current package)")
	fileFlag := flag.String("file", "", "output file name (defaults to current file with 'polygen' suffix)")

	flag.Parse()

	if *typeFlag == "" {
		return nil, fmt.Errorf("type flag is required")
	}
	if *interfaceFlag == "" {
		return nil, fmt.Errorf("interface flag is required")
	}
	if *typesFlag == "" {
		return nil, fmt.Errorf("types flag is required")
	}

	types := make([]TypeMapping, 0)
	for _, t := range strings.Split(*typesFlag, ",") {
		parts := strings.Split(t, "|")
		subType := parts[0]
		typeName := subType
		isPointer := strings.HasPrefix(subType, "*")
		if isPointer {
			subType = subType[1:] // Remove the * prefix
		}
		if len(parts) > 1 {
			typeName = parts[1]
		}
		types = append(types, TypeMapping{
			SubType:   subType,
			TypeName:  typeName,
			IsPointer: isPointer,
		})
	}

	// Set default package name if not provided
	pkg := *packageFlag
	if pkg == "" {
		// Try to detect package from the caller's file (using GOFILE env var set by go:generate)
		if goFile := os.Getenv("GOFILE"); goFile != "" {
			// Read the file and detect its package
			content, err := os.ReadFile(goFile)
			if err == nil {
				// Simple package detection - look for "package" followed by identifier
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "package ") {
						pkg = strings.TrimSpace(strings.TrimPrefix(line, "package "))
						break
					}
				}
			}
		}
		// Fallback to "main" if we couldn't detect the package
		if pkg == "" {
			pkg = "main"
		}
	}

	// Set default file name if not provided
	file := *fileFlag
	if file == "" {
		// If the file flag is empty, use the base name of the first argument
		// and insert 'polygen' before the .go extension
		if len(os.Args) > 0 {
			base := filepath.Base(os.Args[0])
			if strings.HasSuffix(base, ".go") {
				file = strings.TrimSuffix(base, ".go") + "_polygen.go"
			} else {
				file = base + "_polygen.go"
			}
		} else {
			file = "polygen.go"
		}
	}

	return &Config{
		Type:       *typeFlag,
		Interface:  *interfaceFlag,
		Types:      types,
		Descriptor: *descriptorFlag,
		Strict:     *strictFlag,
		Package:    pkg,
		File:       file,
	}, nil
}
