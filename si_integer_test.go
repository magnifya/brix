package humanize

import "testing"

func TestSIWithDigitsIntegerZeros(t *testing.T) {
	tests := []struct {
		name      string
		num       float64
		digits    int
		formatted string
	}{
		{"20", 20, 0, "20 U"},
		{"200", 200, 0, "200 U"},
		{"2000", 2000, 0, "2 kU"},
		{"2.234e-12, 3", 2.234e-12, 3, "2.234 pU"},
		{"2.234e-12, 2", 2.234e-12, 2, "2.23 pU"},
	}
	for _, test := range tests {
		got := SIWithDigits(test.num, test.digits, "U")
		if got != test.formatted {
			t.Errorf("On %v (%v), got %v, wanted %v",
				test.name, test.num, got, test.formatted)
		}
	}
}
