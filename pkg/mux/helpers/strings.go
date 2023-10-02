package helpers

func TruncateFromSecondOccurrence(s string, char rune) string {
	count := 0
	for i, c := range s {
		if c == char {
			count++
			if count == 2 {
				return s[:i]
			}
		}
	}
	// If there isn't a second occurrence of the character, return the original string
	return s
}
