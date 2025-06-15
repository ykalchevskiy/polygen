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

Add a go:generate comment to your source file:

```go
//go:generate polygen -type=ItemValue -interface=IsItemValue -types=ItemValue1|item-value-1,ItemValue2|item-value-2 -descriptor=type -package=domain -file=item_polygen.go
```

### Parameters

- `-type` (required): The name of the polymorphic structure
- `-interface` (required): The name of the interface all subtypes should implement
- `-types` (required): Comma-separated list of subtypes and their type names (format: `SubType|type-name`)
  - Subtype can be prefixed with `*` to indicate a pointer type
  - Type name is optional and defaults to the subtype name
- `-descriptor` (optional): Name of the JSON field to distinguish types (default: "type")
- `-strict` (optional): Enable strict JSON unmarshaling, disallowing unknown fields (default: false)
- `-package` (optional): Package name (defaults to current package)
- `-file` (optional): Output file name (defaults to current file with 'polygen' suffix)

## Features

- Type-safe polymorphic structs with JSON serialization
- Customizable type discriminator field name
- Support for both pointer and value types
- Strict mode for JSON unmarshaling
- Support for patching/updating fields without specifying type
- Automatic type preservation when patching existing values

## Example

Given the following code:

```go
package domain

type IsItemValue interface {
    isItemValue()
}

type ItemValue1 struct {
    Value string
}

func (ItemValue1) isItemValue() {}

type ItemValue2 struct {
    Amount int
}

func (ItemValue2) isItemValue() {}

//go:generate polygen -type=ItemValue -interface=IsItemValue -types=ItemValue1|item-value-1,ItemValue2|item-value-2
```

Running `go generate` will create a file with a polymorphic `ItemValue` type that can be marshaled/unmarshaled to/from JSON:

```go
// Creating a new value
var value ItemValue
json.Unmarshal([]byte(`{"type": "item-value-1", "value": "hello"}`), &value)
// value.IsItemValue will be of type ItemValue1 with Value="hello"

// Patching an existing value - type is preserved
json.Unmarshal([]byte(`{"value": "updated"}`), &value)
// value.IsItemValue still ItemValue1 but with Value="updated"

// Changing type
json.Unmarshal([]byte(`{"type": "item-value-2", "amount": 42}`), &value)
// value.IsItemValue now of type ItemValue2 with Amount=42
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
