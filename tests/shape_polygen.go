// Code generated by polygen; DO NOT EDIT.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	_ IsShape = Circle{}
	_ IsShape = Empty{}
	_ IsShape = (*Group)(nil)
	_ IsShape = (*Polygon)(nil)
	_ IsShape = Rectangle{}
)

// _ShapeTypeRegistry maps concrete types to their type names
var _ShapeTypeRegistry = map[reflect.Type]string{
	reflect.TypeOf((*Circle)(nil)).Elem():    "circle",
	reflect.TypeOf((*Empty)(nil)).Elem():     "empty",
	reflect.TypeOf((*Group)(nil)):            "group",
	reflect.TypeOf((*Polygon)(nil)):          "polygon",
	reflect.TypeOf((*Rectangle)(nil)).Elem(): "rectangle",
}

type Shape struct {
	IsShape
}

func (v Shape) MarshalJSON() ([]byte, error) {
	if v.IsShape == nil {
		return []byte("null"), nil
	}

	// Marshal the implementation first to get its fields
	implData, err := json.Marshal(v.IsShape)
	if err != nil {
		return nil, fmt.Errorf("polygen: cannot marshal IsShape for Shape: %v", err)
	}

	if bytes.Equal(implData, []byte("null")) {
		return implData, nil
	}

	typeName, _, err := _ShapeGetType(v.IsShape)
	if err != nil {
		return nil, fmt.Errorf("polygen: cannot get subtype to marshal for Shape: %v", err)
	}

	// If it's an empty object, just return descriptor
	if string(implData) == "{}" {
		return []byte(fmt.Sprintf("{\"%s\":\"%s\"}", "type", typeName)), nil
	}

	// Otherwise, combine descriptor with implementation fields
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("{\"%s\":\"%s\",", "type", typeName))
	buf.Write(implData[1:])

	return buf.Bytes(), nil
}

func (v *Shape) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = Shape{}
		return nil
	}

	var currTypeName string
	var currTypeAsPointer bool
	if v.IsShape != nil {
		var err error
		currTypeName, currTypeAsPointer, err = _ShapeGetType(v.IsShape)
		if err != nil {
			return fmt.Errorf("polygen: cannot get subtype to unmarshal for Shape: %v", err)
		}
	}
	_ = currTypeAsPointer // In case of all subtypes being pointers, we must just ignore this

	// First decode just the type field
	typeData := struct {
		Type string `json:"type"`
	}{
		Type: currTypeName,
	}
	if err := json.Unmarshal(data, &typeData); err != nil {
		return fmt.Errorf("polygen: cannot unmarshal descriptor field type for Shape: %v", err)
	}

	if typeData.Type == "" {
		return fmt.Errorf("polygen: missing descriptor field type for Shape")
	}

	typeName := typeData.Type

	var value IsShape
	switch typeName {
	case "circle":
		if currTypeName == "circle" {
			if currTypeAsPointer {
				vv := v.IsShape.(*Circle)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Circle for Shape: %v", err)
				}
				value = vv
			} else {
				vv := v.IsShape.(Circle)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Circle for Shape: %v", err)
				}
				value = vv
			}
		} else {
			var vv Circle
			if err := json.Unmarshal(data, &vv); err != nil {
				return fmt.Errorf("polygen: cannot unmarshal Circle for Shape: %v", err)
			}
			value = vv
		}
	case "empty":
		if currTypeName == "empty" {
			if currTypeAsPointer {
				vv := v.IsShape.(*Empty)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Empty for Shape: %v", err)
				}
				value = vv
			} else {
				vv := v.IsShape.(Empty)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Empty for Shape: %v", err)
				}
				value = vv
			}
		} else {
			var vv Empty
			if err := json.Unmarshal(data, &vv); err != nil {
				return fmt.Errorf("polygen: cannot unmarshal Empty for Shape: %v", err)
			}
			value = vv
		}
	case "group":
		var vv *Group
		if currTypeName == "group" {
			vv = v.IsShape.(*Group)
		}
		if err := json.Unmarshal(data, &vv); err != nil {
			return fmt.Errorf("polygen: cannot unmarshal Group for Shape: %v", err)
		}
		value = vv
	case "polygon":
		var vv *Polygon
		if currTypeName == "polygon" {
			vv = v.IsShape.(*Polygon)
		}
		if err := json.Unmarshal(data, &vv); err != nil {
			return fmt.Errorf("polygen: cannot unmarshal Polygon for Shape: %v", err)
		}
		value = vv
	case "rectangle":
		if currTypeName == "rectangle" {
			if currTypeAsPointer {
				vv := v.IsShape.(*Rectangle)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Rectangle for Shape: %v", err)
				}
				value = vv
			} else {
				vv := v.IsShape.(Rectangle)
				if err := json.Unmarshal(data, &vv); err != nil {
					return fmt.Errorf("polygen: cannot unmarshal Rectangle for Shape: %v", err)
				}
				value = vv
			}
		} else {
			var vv Rectangle
			if err := json.Unmarshal(data, &vv); err != nil {
				return fmt.Errorf("polygen: cannot unmarshal Rectangle for Shape: %v", err)
			}
			value = vv
		}
	default:
		return fmt.Errorf("polygen: unknown subtype for Shape: %s", typeName)
	}

	*v = Shape{
		IsShape: value,
	}
	return nil
}

func _ShapeGetType(v IsShape) (name string, asPointer bool, _ error) {
	t := reflect.TypeOf(v)
	typeName, ok := _ShapeTypeRegistry[t]
	if ok {
		return typeName, false, nil
	}
	// A pointer can be manually used for a value type as it also implements the interface
	if t.Kind() == reflect.Ptr {
		typeName, ok = _ShapeTypeRegistry[t.Elem()]
		if ok {
			return typeName, true, nil
		}
	}
	return "", false, fmt.Errorf("unknown subtype: %v", t)
}
