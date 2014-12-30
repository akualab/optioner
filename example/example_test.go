package example

import "testing"

func TestExample(t *testing.T) {

	myFunc := func(n int) int { return 2 * n }
	ex := NewExample("test", N(22), Ff(myFunc))

	if ex.N != 22 {
		t.Errorf("N is %d, expected 22", ex.N)
	}

	if ex.ff(10) != 20 {
		t.Errorf("ff(10) is %d, expected 20", ex.ff(10))
	}

	if cap(ex.FSlice) != 100 {
		t.Errorf("FSlice cap is %d, expected 100", cap(ex.FSlice))
	}

	// Change one of the options. prev has the previous value.
	prev := ex.Option(N(5))

	if ex.N != 5 {
		t.Errorf("N is %d, expected 22", ex.N)
	}

	// restore previous value.
	ex.Option(prev)

	if ex.N != 22 {
		t.Errorf("N is %d, expected 22", ex.N)
	}

}

// All options must be properly rollbacked (i.e. reverted to old values).
func TestExampleRollback(t *testing.T) {

	ex := NewExample("test", N(22), Map(map[string]int{"one": 1}))

	if ex.N != 22 {
		t.Errorf("N is %d, expected 22", ex.N)
	}

	if len(ex.Map) != 1 {
		t.Errorf("len(ex.Map) is %d, expected 1", ex.ff(10))
	}

	// Change two options. original keeps the previous values.
	original := ex.Option(N(33), Map(map[string]int{"one": 1, "two": 2}))

	if ex.N != 33 {
		t.Errorf("N is %d, expected 33", ex.N)
	}

	if len(ex.Map) != 2 {
		t.Errorf("len(ex.Map) is %d, expected 2", ex.ff(10))
	}

	// Restore original value.
	ex.Option(original)

	if ex.N != 22 {
		t.Errorf("N is %d, expected 22", ex.N)
	}

	if len(ex.Map) != 1 {
		t.Errorf("len(ex.Map) is %d, expected 1", ex.ff(10))
	}

}
