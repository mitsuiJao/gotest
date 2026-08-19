package main

import "testing"

/*
func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errof("Add(2, 3) = %d; want %d", got ,want)
	}
}
*/

func TestAddTable(t *testing.T) {
	tests := []struct {
		name 		string
		a, b		int
		expected	int
	}{
		{"positive", 2, 3, 5},
		{"negative", -1, -1, -2},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.expected{
				t.Errorf("got %d; want %d", got, tt.expected)
			}
		})
	}
}
