package main

// Config represents the internal configuration used by the generator
type Config struct {
	Type       string
	Interface  string
	Descriptor string
	Package    string
	Types      []TypeMapping
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
	Types []TypeConfig `json:"types"`
	// StrictByDefault determines if strict mode should be enabled by default for all types
	StrictByDefault bool `json:"strictByDefault,omitempty"`
	// PointerByDefault determines if pointer mode should be enabled by default for all subtypes
	PointerByDefault bool `json:"pointerByDefault,omitempty"`
	// DefaultDescriptor is the default JSON field name to distinguish types
	DefaultDescriptor string `json:"defaultDescriptor,omitempty"`
}

// SubtypeConfig represents configuration for a subtype
type SubtypeConfig struct {
	// Name is the JSON type name, defaults to the subtype name in snake_case if not specified
	Name *string `json:"name,omitempty"`
	// Pointer indicates if this type should be used as a pointer
	Pointer *bool `json:"pointer,omitempty"`
}

// TypeConfig represents configuration for a single polymorphic type
type TypeConfig struct {
	// Type is the name of the polymorphic structure to generate
	Type string `json:"type"`
	// Interface is the name of the interface all subtypes implement
	Interface string `json:"interface"`
	// Package is the package name for the generated file
	Package string `json:"package"`
	// Descriptor is the JSON field name to distinguish types
	Descriptor string `json:"descriptor,omitempty"`
	// Directory is the output directory path, relative to the config file
	Directory string `json:"directory,omitempty"`
	// Filename is the output filename, defaults to <type>_polygen.go in snake_case
	Filename string `json:"filename,omitempty"`
	// Strict enables strict JSON unmarshaling for this type
	Strict *bool `json:"strict,omitempty"`
	// Subtypes maps Go types to their configurations
	Subtypes map[string]SubtypeConfig `json:"subtypes"`
}
