# polygen

[![Go Reference](https://pkg.go.dev/badge/github.com/ykalchevskiy/polygen.svg)](https://pkg.go.dev/github.com/ykalchevskiy/polygen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ykalchevskiy/polygen)](https://golang.org/dl/)

A Go code generator for creating polymorphic data structures with JSON serialization support.

The tool generates Go code for handling polymorphic types in JSON,
allowing different concrete types to be serialized and
deserialized based on a type discriminator field.

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
    "$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/refs/tags/v1.1.0/schema.json",
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
- Global and per-type discriminator field name
- Global and per-type strict mode settings
- Global and per-subtype pointer mode settings
- Global and per-type build tag constraints
- Custom output paths relative to config file
- Global and per-type generation of a per-subtype `<Name>() string` method

### Schema

Per-type settings override their corresponding global defaults. If neither is set,
polygen uses the defaults listed below. The JSON schema is available in [schema.json](./schema.json)
and can be referenced from `.polygen.json` for editor validation and completion.

- `strictByDefault` (optional): Enable strict mode by default (does not apply to jsonv2)
- `pointerByDefault` (optional): Mark all subtypes as pointer mode by default
- `defaultDiscriminator` (optional): Default JSON field name for type discrimination (default: "type")
- `defaultBuildTag` (optional): Build constraint for all generated code (e.g., "linux" or "linux && amd64")
- `defaultTypeNameMethod` (optional): Name for a generated subtype `<Name>() string` method returning the type name
- `defaultJsonVersion` (optional): Default JSON library version to target (options: `v1` (default), `v2experiment`, `v2`)

    * `v2experiment` generates code for Go 1.25 and 1.26 and requires building with `GOEXPERIMENT=jsonv2`.

    * `v2` generates code for Go 1.27 and later.

- `types` (required): Array of type configurations:
  - `type` (required): Name of the polymorphic structure
  - `interface` (required): Name of the interface all subtypes implement
  - `package` (required): Package name for generated code
  - `discriminator` (optional): Override default JSON field name
  - `directory` (optional): Output directory path relative to config file
  - `filename` (optional): Output filename (defaults to `<type_in_snake_case>_polygen.go`)
  - `strict` (optional): Override strict mode for this type (does not apply to jsonv2)
  - `defaultSubtype` (optional): Go subtype key from `subtypes` to use when the discriminator field is missing
  - `buildTag` (optional): Override `defaultBuildTag` for this type
  - `jsonVersion` (optional): JSON library version to target for this type (options: `v1` (default), `v2experiment`, `v2`)
  - `typeNameMethod` (optional): Override `defaultTypeNameMethod` for this type; set to `""` to explicitly disable
  - `subtypes` (required): Map of Go type names to their configurations:
    - `name` (optional): JSON discriminator value (defaults to the Go subtype name converted to kebab-case)
    - `pointer` (optional): Use pointer for this type (defaults to `pointerByDefault`)

When `defaultSubtype` is set, its value is the Go subtype key in `subtypes`, not the subtype's JSON `name`.
