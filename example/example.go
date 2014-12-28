package example

import (
	"fmt"
)

// Example is the struct that will hold optional values.
//go:generate optioner -type Example
type Example struct {
	N      int
	FSlice []float64 `json:"float_slice"`
	Map    map[string]int
	Name   string `opt:"-" json:"name"`
	ff     func(int) int
}

// NewExample creates an example.
// name is required.
func NewExample(name string, options ...Option) *Example {

	// Set required values and initialize optional fields with default values.
	ex := &Example{
		Name:   name,
		N:      10,
		FSlice: make([]float64, 0, 100),
		Map:    make(map[string]int),
		ff:     func(n int) int { return n },
	}

	// Set options.
	ex.init(options...)

	fmt.Printf("Example initalized: %+v\n", ex)
	return ex
}
