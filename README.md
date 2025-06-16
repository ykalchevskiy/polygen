# polygen

[![Go Reference](https://pkg.go.dev/badge/github.com/ykalchevskiy/polygen.svg)](https://pkg.go.dev/github.com/ykalchevskiy/polygen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ykalchevskiy/polygen)](https://goreportcard.com/report/github.com/ykalchevskiy/polygen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ykalchevskiy/polygen)](https://golang.org/dl/)

A Go code generator for creating polymorphic data structures with JSON serialization support.

## Installation

```bash
go install github.com/ykalchevskiy/polygen@latest
```

## Usage

Create a `.polygen.json` file in your project root:

```json
{
    "$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/main/schema.json",
    "strictByDefault": true,
    "defaultDescriptor": "kind",
    "types": [
        {
            "type": "Item",
            "interface": "IsItem",
            "package": "main",
            "directory": "pkg",
            "subtypes": {
                "TextItem": {
                    "name": "text"
                },
                "ImageItem": {
                    "name": "image",
                    "pointer": true
                }
            }
        }
    ]
}
```

Then run:

```bash
polygen
```

## Configuration

The JSON configuration file supports:

- Multiple type definitions in a single file
- Global and per-type strict mode settings
- Global and per-subtype strict mode settings
- Global and per-type descriptor field name
- Simpler subtype configuration with pointer settings
- Custom output paths relative to config file
- Default kebab-case type names for subtypes

### Schema

The configuration follows this structure:

- `strictByDefault` (optional): Enable strict mode by default
- `pointerByDefault` (optional): Enable pointer mode by default
- `defaultDescriptor` (optional): Default JSON field name for type discrimination (default: "type")
- `types` (required): Array of type configurations:
  - `type` (required): Name of the polymorphic structure
  - `interface` (required): Name of the interface all subtypes implement
  - `package` (required): Package name for generated code
  - `descriptor` (optional): Override default JSON field name
  - `directory` (optional): Output directory path relative to config file
  - `filename` (optional): Output filename (defaults to <type>_polygen.go)
  - `strict` (optional): Override strict mode for this type
  - `subtypes` (required): Map of Go type names to their configurations:
    - `name` (optional): JSON type name (defaults to subtype name in kebab-case)
    - `pointer` (optional): Use pointer for this type (defaults to `pointerByDefault`)

## Features

- Type-safe polymorphic structs with JSON serialization
- Customizable type discriminator field name
- Support for both pointer and value types
- Strict mode for JSON unmarshaling
- Support for patching/updating fields without specifying type
- Automatic type preservation when patching existing values
- Default kebab-case type names for cleaner JSON

## Example

Given the following Go code:

```go
type IsItem interface {
    isItem()
}

type TextItem struct {
    Content string
}

func (TextItem) isItem() {}

type ImageItem struct {
    URL    string
    Width  int
    Height int
}

func (ImageItem) isItem() {}
```

The generated code allows you to marshal/unmarshal your types to/from JSON:

```go
// Creating new values
items := []Item{
    {IsItem: TextItem{Content: "Hello, World!"}},
    {IsItem: ImageItem{URL: "https://example.com/image.jpg"}},
}

// Marshaling to JSON
// {"kind": "text", "content": "Hello, World!"}
// {"kind": "image", "url": "https://example.com/image.jpg"}

// Unmarshaling with type changes
var item Item
json.Unmarshal([]byte(`{"kind": "text", "content": "hello"}`), &item)
json.Unmarshal([]byte(`{"content": "updated"}`), &item)  // Updates just content
json.Unmarshal([]byte(`{"kind": "image", "url": "pic.jpg"}`), &item)  // Changes type
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
