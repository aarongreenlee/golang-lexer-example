// Package
package main

import (
	"fmt"
	"os"
)

func main() {
	cases := []string{
		"when we have a lot of words we want to preserve the white space",
		"  which   is  harder to   do for     some examples    than others",
	}

	for i, t := range cases {
		r, err := Parse(t)

		if err != nil {
			fmt.Printf(`Unable to parse the input %d %q\n\tError: %s\n`, i+1, r, err)
			os.Exit(2)
		}

		fmt.Printf(`
Procesing Input %d
	Given    %q
	Parsed   %q
`, i+1, t, r)
	}
}

type itemKind int

const (
	itemWhiteSpace itemKind = iota
	itemEOF
	itemText
)

// Parse receives text and returns a string where the characters in a word have
// ben reversed; however, any whitespace was preserved (e.g., you can not simply
// reverse the entire string provided).
func Parse(text string) (string, error) {
	return text, nil
}
