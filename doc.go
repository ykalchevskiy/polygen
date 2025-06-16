/*
Package polygen is a code generator for creating polymorphic data structures with JSON serialization support.

The tool generates Go code for handling polymorphic types in JSON, allowing different concrete types
to be serialized and deserialized based on a type discriminator field.

Usage:

Create a .polygen.json configuration file:

	{
		"$schema": "https://raw.githubusercontent.com/ykalchevskiy/polygen/main/schema.json",
		"strictByDefault": true,
		"defaultDescriptor": "kind",
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

Configuration options:

	strictByDefault    Enable strict mode by default for all types (optional)
	pointerByDefault   Enable pointer mode by default for all subtypes (optional)
	defaultDescriptor  Default JSON field name for type discrimination (default: "type")
	types             Array of type configurations with the following fields:
	  - type          Name of the polymorphic structure
	  - interface     Name of the interface all subtypes implement
	  - package       Package name for the generated file
	  - descriptor    Override default descriptor field name (optional)
	  - directory     Output directory path relative to config file (optional)
	  - filename      Output filename (defaults to <type>_polygen.go)
	  - strict        Override strict mode for this type (optional)
	  - subtypes      Map of Go types to their configurations:
	                  - name: JSON type name (optional, defaults to subtype in kebab-case)
	                  - pointer: Use pointer for this type (optional, default: false)

Example:

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

This will generate a polymorphic Item type that can handle both TextItem and ImageItem:

	{"kind": "text", "content": "hello"}     // TextItem
	{"kind": "image", "width": 800, ...}     // ImageItem

The generated code also supports patching/merging fields of existing values without changing their type:

	var item Item
	json.Unmarshal([]byte(`{"kind": "text", "content": "hello"}`), &item)
	json.Unmarshal([]byte(`{"content": "updated"}`), &item)  // Updates just the content field
	json.Unmarshal([]byte(`{"kind": "image", "width": 800}`), &item)  // Changes type to ImageItem
*/
package main
