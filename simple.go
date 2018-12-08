// Package
package parserlexer

// ParseSimple is a simple routine to preserve whitespace while reversing the
// characters in words.
func SimpleParser(input string) string {
	var result string
	var word string
	for _, char := range input {
		c := string(char)
		if c == " " {
			// Clean-up the accumulated word
			if len(word) > 0 {
				result += reverse(word)
			}
			result += " "
			continue
		}
	}

	if len(word) > 0 {
		result += reverse(word)
	}

	return result
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
