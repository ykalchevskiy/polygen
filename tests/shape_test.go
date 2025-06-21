package tests

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var marshalTests = []struct {
	name    string
	shape   Shape
	want    string
	wantErr bool
}{
	{
		name:  "circle",
		shape: Shape{IsShape: Circle{Radius: 5.0}},
		want:  `{"type":"circle","Radius":5}`,
	},
	{
		name: "rectangle with style",
		shape: Shape{IsShape: Rectangle{
			Width:  10,
			Height: 20,
			Style: struct {
				Color string
				Fill  bool
			}{
				Color: "red",
				Fill:  true,
			},
		}},
		want: `{"type":"rectangle","Width":10,"Height":20,"Style":{"Color":"red","Fill":true}}`,
	},
	{
		name: "polygon with points",
		shape: Shape{IsShape: &Polygon{
			Points: []struct {
				X float64
				Y float64
			}{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
				{X: 0, Y: 1},
			},
			Labels: []string{"A", "B", "C"},
		}},
		want: `{"type":"polygon","Points":[{"X":0,"Y":0},{"X":1,"Y":1},{"X":0,"Y":1}],"Labels":["A","B","C"]}`,
	},
	{
		name: "group with attributes",
		shape: Shape{IsShape: &Group{
			Name: "test",
			Attributes: map[string]interface{}{
				"visible": true,
				"layer":   1,
				"tags":    []string{"test", "example"},
			},
		}},
		want: `{"type":"group","Name":"test","Attributes":{"layer":1,"tags":["test","example"],"visible":true}}`,
	},
	{
		name:  "nil value",
		shape: Shape{},
		want:  "null",
	},
	{
		name:  "nil interface",
		shape: Shape{IsShape: (*Group)(nil)},
		want:  "null",
	},
	{
		name:  "empty type",
		shape: Shape{IsShape: Empty{}},
		want:  `{"type":"empty"}`,
	},
}

func runMarshalTests(t *testing.T, tests []struct {
	name    string
	shape   Shape
	want    string
	wantErr bool
}, isStrict bool) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got []byte
				err error
			)

			if isStrict {
				strict := ShapeStrict{IsShape: tt.shape.IsShape}
				got, err = json.Marshal(strict)
			} else {
				got, err = json.Marshal(tt.shape)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("%s.MarshalJSON() error = %v, wantErr %v", t.Name(), err, tt.wantErr)
				return
			}
			if err == nil {
				var gotObj, wantObj interface{}
				if err := json.Unmarshal(got, &gotObj); err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}
				if err := json.Unmarshal([]byte(tt.want), &wantObj); err != nil {
					t.Errorf("Failed to unmarshal expected: %v", err)
					return
				}
				if !reflect.DeepEqual(gotObj, wantObj) {
					t.Errorf("%s.MarshalJSON() = %v, want %v", t.Name(), string(got), tt.want)
				}
			}
		})
	}
}

func TestShapeMarshalJSON(t *testing.T) {
	t.Run("non-strict", func(t *testing.T) {
		runMarshalTests(t, marshalTests, false)
	})

	t.Run("strict", func(t *testing.T) {
		runMarshalTests(t, marshalTests, true)
	})
}

type testCase struct {
	name          string
	json          string
	want          interface{}
	wantErr       bool
	strictOnly    bool // test applies to strict mode only
	nonStrictOnly bool // test applies to non-strict mode only
}

var unmarshalTests = []testCase{
	{
		name: "circle",
		json: `{"type":"circle","Radius":5}`,
		want: Shape{IsShape: Circle{Radius: 5.0}},
	},
	{
		name: "rectangle",
		json: `{"type":"rectangle","Width":10,"Height":20,"Style":{"Color":"red","Fill":true}}`,
		want: Shape{IsShape: Rectangle{
			Width:  10,
			Height: 20,
			Style: struct {
				Color string
				Fill  bool
			}{
				Color: "red",
				Fill:  true,
			},
		}},
	},
	{
		name: "polygon",
		json: `{"type":"polygon","Points":[{"X":0,"Y":0},{"X":1,"Y":1}],"Labels":["A","B"]}`,
		want: Shape{IsShape: &Polygon{
			Points: []struct {
				X float64
				Y float64
			}{
				{X: 0, Y: 0},
				{X: 1, Y: 1},
			},
			Labels: []string{"A", "B"},
		}},
	},
	{
		name: "group",
		json: `{"type":"group","Name":"test","Attributes":{"visible":true,"layer":1}}`,
		want: Shape{IsShape: &Group{
			Name: "test",
			Attributes: map[string]interface{}{
				"visible": true,
				"layer":   float64(1),
			},
		}},
	},
	{
		name:    "unknown type",
		json:    `{"type":"unknown"}`,
		wantErr: true,
	},
	{
		name:    "missing type",
		json:    `{"Radius":5}`,
		wantErr: true,
	},
	{
		name:    "invalid json",
		json:    `{`,
		wantErr: true,
	},
	{
		name: "null",
		json: "null",
		want: Shape{},
	},
	{
		name: "empty type",
		json: `{"type":"empty"}`,
		want: Shape{IsShape: Empty{}},
	},
}

// Test cases that are specific to strict or non-strict mode
var extraUnmarshalTests = []testCase{
	{
		name:          "with extra fields",
		json:          `{"type":"circle","Radius":5,"extra":"field"}`,
		want:          Shape{IsShape: Circle{Radius: 5.0}}, // accepted in non-strict mode
		nonStrictOnly: true,
	},
	{
		name:       "with extra fields",
		json:       `{"type":"circle","Radius":5,"extra":"field"}`,
		wantErr:    true, // error in strict mode
		strictOnly: true,
	},
}

// runUnmarshalTests runs the given test cases for both Shape and ShapeStrict types
func runUnmarshalTests(t *testing.T, tests []testCase, isStrict bool) {
	for _, tt := range tests {
		// Skip tests that are meant for the other mode
		if (tt.strictOnly && !isStrict) || (tt.nonStrictOnly && isStrict) {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			var err error
			if isStrict {
				var got ShapeStrict
				err = json.Unmarshal([]byte(tt.json), &got)
				if (err != nil) != tt.wantErr {
					t.Errorf("ShapeStrict.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil && !reflect.DeepEqual(Shape(got), tt.want) {
					t.Errorf("ShapeStrict.UnmarshalJSON() = %+v, want %+v", got, tt.want)
				}
			} else {
				var got Shape
				err = json.Unmarshal([]byte(tt.json), &got)
				if (err != nil) != tt.wantErr {
					t.Errorf("Shape.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if err == nil && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Shape.UnmarshalJSON() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestShapeUnmarshalJSON(t *testing.T) {
	t.Run("non-strict", func(t *testing.T) {
		runUnmarshalTests(t, append(unmarshalTests, extraUnmarshalTests...), false)
	})

	t.Run("strict", func(t *testing.T) {
		runUnmarshalTests(t, append(unmarshalTests, extraUnmarshalTests...), true)
	})
}

func TestShapeMarshalJSONIndent(t *testing.T) {
	shape := Shape{IsShape: &Polygon{
		Points: []struct {
			X float64
			Y float64
		}{
			{X: 0, Y: 0},
			{X: 1, Y: 1},
			{X: 0, Y: 1},
		},
		Labels: []string{"A", "B", "C"},
	}}

	got, err := json.MarshalIndent(shape, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent failed: %v", err)
	}

	want := `{
  "type": "polygon",
  "Points": [
    {
      "X": 0,
      "Y": 0
    },
    {
      "X": 1,
      "Y": 1
    },
    {
      "X": 0,
      "Y": 1
    }
  ],
  "Labels": [
    "A",
    "B",
    "C"
  ]
}`

	if string(got) != want {
		t.Errorf("MarshalIndent produced incorrect output.\nGot:\n%s\nWant:\n%s", string(got), want)
	}
}

func TestShapeJSONStreamDecoder(t *testing.T) {
	t.Run("non-strict", func(t *testing.T) { runJSONStreamDecoder(t, false) })
	t.Run("strict", func(t *testing.T) { runJSONStreamDecoder(t, true) })
}

func runJSONStreamDecoder(t *testing.T, isStrict bool) {
	input := strings.Join([]string{
		`{"type":"circle","Radius":1}`,
		`{"type":"rectangle","Width":2,"Height":3}`,
		`{"type":"polygon","Points":[{"X":0,"Y":0}]}`,
		`{"type":"group","Name":"test"}`,
	}, "\n")

	decoder := json.NewDecoder(bytes.NewReader([]byte(input)))

	expected := []Shape{
		{IsShape: Circle{Radius: 1}},
		{IsShape: Rectangle{Width: 2, Height: 3}},
		{IsShape: &Polygon{Points: []struct{ X, Y float64 }{{X: 0, Y: 0}}}},
		{IsShape: &Group{Name: "test"}},
	}

	for i, want := range expected {
		if isStrict {
			var got ShapeStrict
			if err := decoder.Decode(&got); err != nil {
				t.Fatalf("Failed to decode item %d: %v", i, err)
			}
			if !reflect.DeepEqual(Shape(got), want) {
				t.Errorf("Item %d: got %+v, want %+v", i, got, want)
			}
		} else {
			var got Shape
			if err := decoder.Decode(&got); err != nil {
				t.Fatalf("Failed to decode item %d: %v", i, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Item %d: got %+v, want %+v", i, got, want)
			}
		}
	}

	// Should be at EOF
	if decoder.More() {
		t.Error("Expected to be at EOF, but decoder.More() returned true")
	}
}

func TestShapeUpdate(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		update  string
		want    Shape
		wantErr bool
	}{
		{
			name:    "update value type (circle)",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"circle","Radius":10}`,
			want:    Shape{IsShape: Circle{Radius: 10.0}},
		},
		{
			name:    "update pointer type (group)",
			initial: `{"type":"group","Name":"test","Attributes":{"active":true}}`,
			update:  `{"type":"group","Name":"updated","Attributes":{"status":"ready"}}`,
			want: Shape{IsShape: &Group{
				Name: "updated",
				Attributes: map[string]interface{}{
					"active": true,
					"status": "ready",
				},
			}},
		},
		{
			name:    "update pointer type (polygon)",
			initial: `{"type":"polygon","Points":[{"X":0,"Y":0}],"Labels":["A"]}`,
			update:  `{"type":"polygon","Points":[{"X":1,"Y":1}],"Labels":["B"]}`,
			want: Shape{IsShape: &Polygon{
				Points: []struct{ X, Y float64 }{{X: 1, Y: 1}},
				Labels: []string{"B"},
			}},
		},
		{
			name:    "update type field from value to value",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"rectangle", "Width":10,"Height":20}`,
			want: Shape{IsShape: Rectangle{
				Width:  10.0,
				Height: 20.0,
			}},
		},
		{
			name:    "update type field from value to value empty",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"rectangle"}`,
			want:    Shape{IsShape: Rectangle{}},
		},
		{
			name:    "update type field from value to pointer",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"polygon", "Points":[{"X":0,"Y":0}],"Labels":["A"]}`,
			want: Shape{IsShape: &Polygon{
				Points: []struct {
					X float64
					Y float64
				}{{X: 0, Y: 0}},
				Labels: []string{"A"}}},
		},
		{
			name:    "update type field from value to pointer empty",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"polygon"}`,
			want:    Shape{IsShape: &Polygon{}},
		},
		{
			name:    "update type field from pointer to value",
			initial: `{"type":"group","Name":"test","Attributes":{"active":true}}`,
			update:  `{"type":"circle","Radius":10}`,
			want:    Shape{IsShape: Circle{Radius: 10.0}},
		},
		{
			name:    "update type field from pointer to value empty",
			initial: `{"type":"group","Name":"test","Attributes":{"active":true}}`,
			update:  `{"type":"circle"}`,
			want:    Shape{IsShape: Circle{}},
		},
		{
			name:    "update type field from pointer to pointer",
			initial: `{"type":"polygon","Points":[{"X":0,"Y":0}],"Labels":["A"]}`,
			update:  `{"type":"group","Name":"updated","Attributes":{"status":"ready"}}`,
			want: Shape{IsShape: &Group{
				Name: "updated",
				Attributes: map[string]interface{}{
					"status": "ready",
				},
			}},
		},
		{
			name:    "update type field from pointer to pointer empty",
			initial: `{"type":"polygon","Points":[{"X":0,"Y":0}],"Labels":["A"]}`,
			update:  `{"type":"group"}`,
			want:    Shape{IsShape: &Group{}},
		},
		{
			name:    "update with unknown type",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"unknown"}`,
			wantErr: true,
		},
		{
			name:    "update with invalid JSON",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":"circle","Radius":10,`, // Invalid JSON
			wantErr: true,
		},
		{
			name:    "update with empty JSON value",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{}`,                                // Empty JSON object
			want:    Shape{IsShape: Circle{Radius: 5.0}}, // Should keep the original value
		},
		{
			name:    "update with empty JSON pointer",
			initial: `{"type":"group","Name":"updated","Attributes":{"status":"ready"}}`,
			update:  `{}`, // Empty JSON object
			want: Shape{IsShape: &Group{
				Name: "updated",
				Attributes: map[string]interface{}{
					"status": "ready",
				},
			}},
		},
		{
			name:    "update with null",
			initial: `{"type":"circle","Radius":5}`,
			update:  `null`, // Null value
			want:    Shape{},
		},
		{
			name:    "update without type value",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"type":""}`,
			wantErr: true,
		},
		{
			name:    "update without type pointer",
			initial: `{"type":"polygon"}`,
			update:  `{"type":""}`,
			wantErr: true,
		},
		{
			name:    "update without data value",
			initial: `{"type":"circle","Radius":5}`,
			update:  `{"Radius":10}`,
			want:    Shape{IsShape: Circle{Radius: 10.0}}, // Should update the field
		},
		{
			name:    "update without data pointer",
			initial: `{"type":"group","Name":"updated","Attributes":{"status":"ready"}}`,
			update:  `{"Attributes":{"visible":true}}`,
			want: Shape{IsShape: &Group{
				Name: "updated",
				Attributes: map[string]interface{}{
					"status":  "ready",
					"visible": true,
				},
			}},
		},
	}

	for _, mode := range []struct {
		name     string
		isStrict bool
	}{
		{"non-strict", false},
		{"strict", true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					if mode.isStrict {
						var err error

						var got ShapeStrict
						if err = json.Unmarshal([]byte(tt.initial), &got); err != nil {
							t.Fatalf("Failed to unmarshal initial value: %v", err)
						}
						err = json.Unmarshal([]byte(tt.update), &got)
						if (err != nil) != tt.wantErr {
							t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
							return
						}
						if err == nil && !reflect.DeepEqual(got.IsShape, tt.want.IsShape) {
							t.Errorf("UnmarshalJSON() = %+v, want %+v", got.IsShape, tt.want.IsShape)
						}
					} else {
						var err error

						var got Shape
						if err = json.Unmarshal([]byte(tt.initial), &got); err != nil {
							t.Fatalf("Failed to unmarshal initial value: %v", err)
						}
						err = json.Unmarshal([]byte(tt.update), &got)
						if (err != nil) != tt.wantErr {
							t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
							return
						}
						if err == nil && !reflect.DeepEqual(got.IsShape, tt.want.IsShape) {
							t.Errorf("UnmarshalJSON() = %+v, want %+v", got.IsShape, tt.want.IsShape)
						}
					}
				})
			}
		})
	}
}

func TestShapeSettability(t *testing.T) {
	t.Run("value subtype is not settable", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("The code did not panic")
			}
			if !strings.Contains(r.(string), "reflect: reflect.Value.SetInt using unaddressable value") {
				t.Errorf("Expected panic about non-settable value, got: %v", r)
			}
		}()

		var value Shape
		err := json.Unmarshal([]byte(`{"type":"circle","Radius":5}`), &value)
		if err != nil {
			t.Fatalf("Failed to unmarshal value: %v", err)
		}

		reflect.ValueOf(&value).Elem().FieldByName("IsShape").Elem().FieldByName("Radius").SetInt(10) // This should panic

		t.Errorf("Expected panic when setting value field, but did not panic")
	})

	t.Run("pointer subtype is settable", func(t *testing.T) {
		var pointer Shape
		err := json.Unmarshal([]byte(`{"type":"group","Name":"test","Attributes":{"active":true}}`), &pointer)
		if err != nil {
			t.Fatalf("Failed to unmarshal pointer: %v", err)
		}

		reflect.ValueOf(&pointer).Elem().FieldByName("IsShape").Elem().Elem().FieldByName("Name").SetString("updated")

		if pointer.IsShape.(*Group).Name != "updated" {
			t.Errorf("Expected pointer name to be updated, got: %s", pointer.IsShape.(*Group).Name)
		}
	})
}

func TestShape(t *testing.T) {
	t.Run("non_strict/value for value type", func(t *testing.T) {
		shape := Shape{IsShape: Circle{Radius: 5.0}}
		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal value type: %v", err)
		}
		expected := `{"type":"circle","Radius":5}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"circle", "Radius":10}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal value type: %v", err)
		}
		if shape.IsShape.(Circle).Radius != 10.0 {
			t.Errorf("Expected radius to be 10.0, got %f", shape.IsShape.(Circle).Radius)
		}
	})

	t.Run("strict/value for value type", func(t *testing.T) {
		shape := ShapeStrict{IsShape: Circle{Radius: 5.0}}
		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal value type: %v", err)
		}
		expected := `{"type":"circle","Radius":5}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"circle", "Radius":10}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal value type: %v", err)
		}
		if shape.IsShape.(Circle).Radius != 10.0 {
			t.Errorf("Expected radius to be 10.0, got %f", shape.IsShape.(Circle).Radius)
		}
	})

	t.Run("non_strict/pointer for value type", func(t *testing.T) {
		shape := Shape{IsShape: &Circle{Radius: 5.0}}
		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal value type: %v", err)
		}
		expected := `{"type":"circle","Radius":5}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"circle", "Radius":10}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal value type: %v", err)
		}
		if shape.IsShape.(*Circle).Radius != 10.0 {
			t.Errorf("Expected radius to be 10.0, got %f", shape.IsShape.(*Circle).Radius)
		}
	})

	t.Run("strict/pointer for value type", func(t *testing.T) {
		shape := ShapeStrict{IsShape: &Circle{Radius: 5.0}}
		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal value type: %v", err)
		}
		expected := `{"type":"circle","Radius":5}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"circle", "Radius":10}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal value type: %v", err)
		}
		if shape.IsShape.(*Circle).Radius != 10.0 {
			t.Errorf("Expected radius to be 10.0, got %f", shape.IsShape.(*Circle).Radius)
		}
	})

	t.Run("non_strict/pointer for pointer type", func(t *testing.T) {
		shape := Shape{IsShape: &Group{
			Name: "test",
			Attributes: map[string]interface{}{
				"active": true,
			},
		}}

		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal pointer type: %v", err)
		}
		expected := `{"type":"group","Name":"test","Attributes":{"active":true}}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"group", "Name":"updated"}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal pointer type: %v", err)
		}
		if shape.IsShape.(*Group).Name != "updated" {
			t.Errorf("Expected name to be 'updated', got: %s", shape.IsShape.(*Group).Name)
		}
	})

	t.Run("strict/pointer for pointer type", func(t *testing.T) {
		shape := ShapeStrict{IsShape: &Group{
			Name: "test",
			Attributes: map[string]interface{}{
				"active": true,
			},
		}}

		data, err := json.Marshal(shape)
		if err != nil {
			t.Fatalf("Failed to marshal pointer type: %v", err)
		}
		expected := `{"type":"group","Name":"test","Attributes":{"active":true}}`
		if string(data) != expected {
			t.Errorf("Expected %s, got %s", expected, string(data))
		}

		if err := json.Unmarshal([]byte(`{"type":"group", "Name":"updated"}`), &shape); err != nil {
			t.Fatalf("Failed to unmarshal pointer type: %v", err)
		}
		if shape.IsShape.(*Group).Name != "updated" {
			t.Errorf("Expected name to be 'updated', got: %s", shape.IsShape.(*Group).Name)
		}
	})

}
