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
	// DefaultDescriptor is the default JSON field name to distinguish types
	DefaultDescriptor string `json:"defaultDescriptor,omitempty"`
}

// SubtypeConfig represents configuration for a subtype
type SubtypeConfig struct {
	// Name is the JSON type name, defaults to the subtype name in snake_case if not specified
	Name *string `json:"name,omitempty"`
	// Pointer indicates if this type should be used as a pointer
	Pointer bool `json:"pointer,omitempty"`
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

// JSONSchema returns the JSON Schema for the configuration file
func JSONSchema() string {
	return `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"required": ["types"],
	"properties": {
		"strictByDefault": {
			"type": "boolean",
			"description": "Enable strict mode by default for all types"
		},
		"defaultDescriptor": {
			"type": "string",
			"description": "Default JSON field name to distinguish types"
		},
		"types": {
			"type": "array",
			"items": {
				"type": "object",
				"required": ["type", "interface", "package", "subtypes"],
				"properties": {
					"type": {
						"type": "string",
						"description": "Name of the polymorphic structure to generate"
					},
					"interface": {
						"type": "string",
						"description": "Name of the interface all subtypes implement"
					},
					"package": {
						"type": "string",
						"description": "Package name for the generated file"
					},
					"descriptor": {
						"type": "string",
						"description": "JSON field name to distinguish types",
						"default": "type"
					},
						"directory": {
						"type": "string",
						"description": "Output directory path relative to the config file"
					},
					"filename": {
						"type": "string",
						"description": "Output filename, defaults to <type>_polygen.go in snake_case"
					},
					"strict": {
						"type": "boolean",
						"description": "Enable strict JSON unmarshaling for this type"
					},
					"subtypes": {
						"type": "object",
						"description": "Map of Go types to their configurations",
						"additionalProperties": {
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"description": "JSON type name, defaults to the subtype name in snake_case if not specified"
								},
								"pointer": {
									"type": "boolean",
									"description": "Indicates if this type should be used as a pointer"
								}
							},
							"required": []
						}
					}
				}
			},
			"minItems": 1
		}
	}
}`
}
