package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

//go:generate go run github.com/ykalchevskiy/polygen

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

func (*ImageItem) isItem() {}

func main() {
	var item Item
	json.Unmarshal([]byte(`{"kind": "image", "width": 800}`), &item) // Changes type to ImageItem
	fmt.Println(item.IsItem)
	json.Unmarshal([]byte(`{"kind": "text", "content": "hello"}`), &item)
	fmt.Println(item.IsItem)
	json.Unmarshal([]byte(`{"content": "updated"}`), &item) // Updates just the content field
	fmt.Println(item.IsItem)
	json.Unmarshal([]byte(`{"kind": "image", "width": 800}`), &item) // Changes type to ImageItem
	fmt.Println(item.IsItem)
	json.Unmarshal([]byte(`{"width": 900}`), &item) // Changes type to ImageItem
	fmt.Println(item.IsItem)

	// Create and marshal items
	items := []Item{
		{
			IsItem: TextItem{Content: "Hello, World!"},
		},
		{
			IsItem: &ImageItem{
				URL:    "https://example.com/image.jpg",
				Width:  800,
				Height: 600,
			},
		},
	}

	for _, item := range items {
		data, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Marshaled item:\n%s\n\n", data)
	}

	// Unmarshal items
	jsonData := []string{
		`{"kind": "text", "content": "Hello from JSON!"}`,
		`{"kind": "image", "url": "https://example.com/pic.jpg", "width": 1024, "height": 768}`,
	}

	for _, data := range jsonData {
		var item Item
		if err := json.Unmarshal([]byte(data), &item); err != nil {
			log.Fatal(err)
		}

		switch v := item.IsItem.(type) {
		case TextItem:
			fmt.Printf("Got text item: %s\n", v.Content)
			// patch item with reflect package
			// reflect.ValueOf(item).FieldByName("IsItem").Elem().FieldByName("Content").SetString("Patched content")
			fmt.Printf("Got text item patched: %s\n", item)
		case *ImageItem:
			fmt.Printf("Got image item: %dx%d at %s\n", v.Width, v.Height, v.URL)

			reflect.ValueOf(item).FieldByName("IsItem").Elem().Elem().FieldByName("URL").SetString("another-url.jpg")
			fmt.Printf("Got image item patched: %s\n", item)
		}
	}
}
