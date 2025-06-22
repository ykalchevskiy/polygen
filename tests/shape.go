package tests

//go:generate go run ..

type IsShape interface {
	isShape()
}

type Circle struct {
	Radius float64
}

func (Circle) isShape() {}

type Rectangle struct {
	Width  float64
	Height float64
	Style  struct {
		Color string
		Fill  bool
	}
}

func (Rectangle) isShape() {}

type Polygon struct {
	Points []struct {
		X float64
		Y float64
	}
	Labels []string
}

func (*Polygon) isShape() {}

type Group struct {
	Name       string
	Attributes map[string]any
}

func (*Group) isShape() {}

type Empty struct{}

func (Empty) isShape() {}
