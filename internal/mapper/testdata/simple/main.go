package main

import "fmt"

// Widget is a simple exported struct for testing AST extraction.
type Widget struct {
	Name string
}

// NewWidget creates a new Widget with the given name.
func NewWidget(name string) *Widget {
	return &Widget{Name: name}
}

func main() {
	w := NewWidget("test")
	fmt.Println(w.Name)
}
