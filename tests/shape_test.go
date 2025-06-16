package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Run polygen to generate test types
	genCmd := exec.Command("go", "run", "..")
	genCmd.Dir = "." // Run in tests directory
	if err := genCmd.Run(); err != nil {
		panic("failed to run polygen: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

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
				if tt.want == "null" {
					if string(got) != tt.want {
						t.Errorf("%s.MarshalJSON() = %v, want %v", t.Name(), string(got), tt.want)
					}
				} else {
					// Compare JSON objects to ignore field order differences
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
	{
		name: "update existing value",
		json: `{"Radius":10}`,
		want: Shape{IsShape: Circle{Radius: 10.0}},
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
				if tt.name == "update existing value" {
					// First unmarshal a circle, then update it
					if err = json.Unmarshal([]byte(`{"type":"circle","Radius":5}`), &got); err != nil {
						t.Fatalf("Failed to unmarshal initial value: %v", err)
					}
				}
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
				if tt.name == "update existing value" {
					// First unmarshal a circle, then update it
					if err = json.Unmarshal([]byte(`{"type":"circle","Radius":5}`), &got); err != nil {
						t.Fatalf("Failed to unmarshal initial value: %v", err)
					}
				}
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
