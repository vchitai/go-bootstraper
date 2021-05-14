package utils

import "unicode"

type gap struct {
	PrevEnd   int
	NextStart int
}
type dfa interface {
	check(s rune) *gap
}

func getDFAs() []dfa {
	return []dfa{
		&dashDFA{},
		&titleDFA{},
		&specialDFA{},
	}
}

var _ dfa = &dashDFA{}
var _ dfa = &titleDFA{}
var _ dfa = &specialDFA{}

// for dash
type dashDFA struct {
	current int
	state   int
	gap
}

func (t *dashDFA) check(c rune) *gap {
	defer func() {
		t.current += 1
	}()
	switch t.state {
	case 0:
		if unicode.IsLower(c) {
			t.state = 1
			t.PrevEnd = t.current + 1
			return nil
		}
	case 1:
		if unicode.IsLower(c) {
			t.state = 1
			t.PrevEnd = t.current + 1
			return nil
		}
		if c == '-' {
			t.state = 2
			return nil
		}
	case 2:
		if unicode.IsLower(c) || unicode.IsUpper(c) {
			t.state = 0
			t.NextStart = t.current
			return &t.gap
		}
	}
	t.state = 0
	return nil
}

// for HelloWorld
type titleDFA struct {
	current int
	state   int
}

func (t *titleDFA) check(c rune) *gap {
	defer func() {
		t.current += 1
	}()
	switch t.state {
	case 0:
		if unicode.IsLower(c) {
			t.state = 1
			return nil
		}
	case 1:
		if unicode.IsLower(c) {
			t.state = 1
			return nil
		}
		if unicode.IsUpper(c) {
			t.state = 0
			return &gap{
				PrevEnd:   t.current,
				NextStart: t.current,
			}
		}
	}
	t.state = 0
	return nil
}

// for URLs case
type specialDFA struct {
	current int
	state   int
}

func (t *specialDFA) check(c rune) *gap {
	defer func() {
		t.current += 1
	}()
	switch t.state {
	case 0:
		if unicode.IsUpper(c) {
			t.state = 1
			return nil
		}
	case 1:
		if unicode.IsLower(c) {
			t.state = 0
			return nil
		}
		if unicode.IsUpper(c) {
			t.state = 2
			return nil
		}
	case 2:
		if unicode.IsLower(c) {
			if c != 's' {
				t.state = 0
				return &gap{
					PrevEnd:   t.current - 1,
					NextStart: t.current - 1,
				}
			} else {
				t.state = 1
				return nil
			}
		}
		if unicode.IsUpper(c) {
			t.state = 2
			return nil
		}
	}
	t.state = 0
	return nil
}
