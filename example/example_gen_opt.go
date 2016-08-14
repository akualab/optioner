// generated by optioner -type Example -m Option; DO NOT EDIT

// Please report issues and submit contributions at:
// http://github.com/akualab/optioner
// optioner is a project of AKUALAB INC.

package example

// Option type is used to set options in Example.
type Option func(*Example) Option

// Option method sets the options. Returns previous option for last arg.
func (t *Example) Option(options ...Option) (previous Option) {
	for _, opt := range options {
		previous = opt(t)
	}
	return previous
}

// N sets a value for instances of type Example.
func N(o int) Option {
	return func(t *Example) Option {
		previous := t.N
		t.N = o
		return N(previous)
	}
}

// FSlice sets a value for instances of type Example.
func FSlice(o []float64) Option {
	return func(t *Example) Option {
		previous := t.FSlice
		t.FSlice = o
		return FSlice(previous)
	}
}

// Map sets a value for instances of type Example.
func Map(o map[string]int) Option {
	return func(t *Example) Option {
		previous := t.Map
		t.Map = o
		return Map(previous)
	}
}

// Func sets a value for instances of type Example.
func Func(o func(int) int) Option {
	return func(t *Example) Option {
		previous := t.ff
		t.ff = o
		return Func(previous)
	}
}
