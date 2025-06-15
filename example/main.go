package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func (ImageItem) isItem() {}

func main() {
	// Create and marshal items
	items := []Item{
		{
			IsItem: TextItem{Content: "Hello, World!"},
		},
		{
			IsItem: ImageItem{
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
		case ImageItem:
			fmt.Printf("Got image item: %dx%d at %s\n", v.Width, v.Height, v.URL)
		}
	}
}
