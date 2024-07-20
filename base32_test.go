package id

import "testing"

var base32Input = [...]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56}

func TestToBase32(t *testing.T) {
	want := "041061050R3GG28A1C60T3GF208H44RM2MB1E60S38DHR78Y3WG228H34GJJC9S854N2PB1D5RQK0C9J6CT3ADHQ7R"

	got := toBase32String(base32Input)
	if got != want {
		t.Errorf("toBase32String(...) = %q; want %q", got, want)
	}
}

func TestToBase32Zero(t *testing.T) {
	input := [fullIDLength]byte{}
	want := "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	got := toBase32String(input)
	if got != want {
		t.Errorf("toBase32String(...) = %q; want %q", got, want)
	}
}

func BenchmarkToBase32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = toBase32String(base32Input)
	}
}
