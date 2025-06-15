/*
Package polygen is a code generator for creating polymorphic data structures with JSON serialization support.

The tool generates Go code for handling polymorphic types in JSON, allowing different concrete types
to be serialized and deserialized based on a type discriminator field.

Usage:

Add a go:generate comment to your source file:

	//go:generate polygen -type=ItemValue -interface=IsItemValue -types=ItemValue1|item-value-1,ItemValue2|item-value-2

Parameters:

	-type        (required) The name of the polymorphic structure
	-interface   (required) The name of the interface all subtypes should implement
	-types       (required) Comma-separated list of subtypes and their type names (format: SubType|type-name)
	             Subtype can be prefixed with * to indicate a pointer type
	             Type name is optional and defaults to the subtype name
	-descriptor  (optional) Name of the JSON field to distinguish types (default: "type")
	-strict      (optional) Enable strict JSON unmarshaling, disallowing unknown fields (default: false)
	-package     (optional) Package name (defaults to current package)
	-file        (optional) Output file name (defaults to current file with 'polygen' suffix)

Example:

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

This will generate a polymorphic ItemValue type that can handle both ItemValue1 and ItemValue2:

	{"type": "item-value-1", "value": "hello"}  // ItemValue1
	{"type": "item-value-2", "amount": 42}      // ItemValue2

The generated code also supports patching/merging fields of existing values without changing their type:

	var item ItemValue
	json.Unmarshal([]byte(`{"type": "item-value-1", "value": "hello"}`), &item)
	json.Unmarshal([]byte(`{"value": "updated"}`), &item)  // Updates just the value field
	json.Unmarshal([]byte(`{"type": "item-value-2", "amount": 42}`), &item)  // Changes type to ItemValue2
*/
package main
