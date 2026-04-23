package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"Charmander Bulbasaur PIKACHU", []string{"charmander", "bulbasaur", "pikachu"}},
		{"  leading and trailing  ", []string{"leading", "and", "trailing"}},
		{"  multiple   spaces   between  ", []string{"multiple", "spaces", "between"}},
		{"", []string{}},
		{"SINGLE", []string{"single"}},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) != len(c.expected) {
			t.Errorf("cleanInput(%q): expected %d words, got %d", c.input, len(c.expected), len(result))
			continue
		}
		for i := range result {
			if result[i] != c.expected[i] {
				t.Errorf("cleanInput(%q): expected %q at index %d, got %q", c.input, c.expected[i], i, result[i])
			}
		}
	}
}
