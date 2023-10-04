package helpers

import "testing"

// TestTruncateFromNthOccurrence was written by JetBrains AI 🤯
func TestTruncateFromNthOccurrence(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		char       rune
		occurrence int
		want       string
	}{
		{
			name:       "Occurrence is zero",
			s:          "Hello, world",
			char:       ',',
			occurrence: 0,
			want:       "Hello, world",
		},
		{
			name:       "Occurrence is 1",
			s:          "Hello, world",
			char:       ',',
			occurrence: 1,
			want:       "Hello",
		},
		{
			name:       "Occurrence is more than string length",
			s:          "Hello world",
			char:       ' ',
			occurrence: 2,
			want:       "Hello world",
		},
		{
			name:       "Char is not in string",
			s:          "Hello world",
			char:       ',',
			occurrence: 1,
			want:       "Hello world",
		},
		{
			name:       "Char is part of string multiple times",
			s:          "Hello world, Gophers",
			char:       ' ',
			occurrence: 2,
			want:       "Hello world,",
		},
		{
			name:       "Empty string",
			s:          "",
			char:       'H',
			occurrence: 1,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateFromNthOccurrence(tt.s, tt.char, tt.occurrence); got != tt.want {
				t.Errorf("TruncateFromNthOccurrence() = %v, want %v", got, tt.want)
			}
		})
	}
}
