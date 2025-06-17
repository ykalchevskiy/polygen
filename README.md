# polygen

[![Go Reference](https://pkg.go.dev/badge/github.com/ykalchevskiy/polygen.svg)](https://pkg.go.dev/github.com/ykalchevskiy/polygen)
[![Go Report Card](https://goreportcard.com/badge/github.com/ykalchevskiy/polygen)](https://goreportcard.com/report/github.com/ykalchevskiy/polygen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ykalchevskiy/polygen)](https://golang.org/dl/)

A Go code generator for creating polymorphic data structures with JSON serialization support.

## Example

Given the following Go code:

```go
type IsItem interface {
    isItem()
}

type TextItem struct {
    Content string `json:"content"`
}

func (TextItem) isItem() {}

type ImageItem struct {
    URL    string `json:"url"`
}

func (*ImageItem) isItem() {}
```

Create a `.polygen.json` file in your project root:

```json
{
    "$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/main/schema.json",
    "types": [
        {
            "type": "Item",
            "interface": "IsItem",
            "package": "main",
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

The generated code allows you to marshal/unmarshal your types to/from JSON:

```go
// unmarshaling
var item Item
json.Unmarshal([]byte(`{"type": "text", "content": "hello"}`), &item)
json.Unmarshal([]byte(`{"content": "updated"}`), &item)  // Updates just content
json.Unmarshal([]byte(`{"type": "image", "url": "pic.jpg"}`), &item)  // Changes type


// marshaling
json.Marshal(Item{IsItem: TextItem{Content: "Hello, World!"}})
// {"type": "text", "content": "Hello, World!"}
json.Marshal(Item{IsItem: &ImageItem{URL: "https://example.com/image.jpg"}})
// {"type": "image", "url": "https://example.com/image.jpg"}
```

## Installation

```bash
go install github.com/ykalchevskiy/polygen@latest
```

## Configuration

The JSON configuration file supports:

- Multiple type definitions in a single file
- Global and per-type descriptor field name
- Global and per-type strict mode settings
- Global and per-subtype pointer mode settings
- Custom output paths relative to config file

### Schema

The configuration follows this structure:

- `strictByDefault` (optional): Enable strict mode by default
- `pointerByDefault` (optional): Mark all subtypes as pointer mode by default
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
