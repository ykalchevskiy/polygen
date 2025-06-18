package tests

//go:generate go run ..

// Common interface for testing value and pointer types
type IsShape interface {
	isShape()
}

// Value type with simple fields
type Circle struct {
	Radius float64
}

func (Circle) isShape() {}

// Value type with nested struct
type Rectangle struct {
	Width  float64
	Height float64
	Style  struct {
		Color string
		Fill  bool
	}
}

func (Rectangle) isShape() {}

// Pointer type with slices
type Polygon struct {
	Points []struct {
		X float64
		Y float64
	}
	Labels []string
}

func (*Polygon) isShape() {}

// Pointer type with maps
type Group struct {
	Name       string
	Attributes map[string]interface{}
}

func (*Group) isShape() {}

// Empty type to test descriptor-only marshaling
type Empty struct{}

func (Empty) isShape() {}
