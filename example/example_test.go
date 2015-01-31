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

func TestPerson(t *testing.T) {

	p := NewPerson("111-222-3333", Name("joe"), Age(22))

	if p.Age != 22 {
		t.Errorf("Age is %d, expected 22", p.Age)
	}

	if p.Name != "joe" {
		t.Errorf("Name is %s, expected joe", p.Name)
	}

}
