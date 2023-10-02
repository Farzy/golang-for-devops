package helpers

func TruncateFromNthOccurrence(s string, char rune, occurrence int) string {
	if occurrence < 1 {
		return s
	}

	count := 0
	for i, c := range s {
		if c == char {
			count++
			if count == occurrence {
				return s[:i]
			}
		}
	}
	// If there is no occurrence of the character, return the original string
	return s
}
