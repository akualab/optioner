package example

import (
	"fmt"
)

// Example is the struct that will hold optional values.
//go:generate optioner -type Example -m Option
type Example struct {
	N      int
	FSlice []float64 `json:"float_slice"`
	Map    map[string]int
	Name   string        `opt:"-" json:"name"`
	ff     func(int) int `opt:"Func"`
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
	ex.Option(options...)

	fmt.Printf("Example initalized: %+v\n", ex)
	return ex
}

// Person is a human.
//go:generate optioner -type Person
type Person struct {
	Name string
	Age  int
	ssn  string `opt:"-"`
}

// NewExample creates an example.
// name is required.
func NewPerson(ssn string, options ...optPerson) *Person {

	// Set required values and initialize optional fields with default values.
	p := &Person{
		ssn: ssn,
	}

	// Set options.
	p.Option(options...)

	fmt.Printf("Person initalized: %+v\n", p)
	return p
}
