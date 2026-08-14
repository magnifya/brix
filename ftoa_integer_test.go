package humanize

import "testing"

func TestFtoaWithDigitsIntegerZeros(t *testing.T) {
	testList{
		{"20.0, 0", FtoaWithDigits(20.0, 0), "20"},
		{"200.0, 0", FtoaWithDigits(200.0, 0), "200"},
		{"2000.0, 0", FtoaWithDigits(2000.0, 0), "2000"},
		{"2.0, 0", FtoaWithDigits(2.0, 0), "2"},
		{"1.23, 0", FtoaWithDigits(1.23, 0), "1"},
		{"1.23, 1", FtoaWithDigits(1.23, 1), "1.2"},
		{"1.23, 2", FtoaWithDigits(1.23, 2), "1.23"},
		{"1.230000, 6", FtoaWithDigits(1.230000, 6), "1.23"},
	}.validate(t)
}
