// Package
package main

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
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

const eof = -1

const (
	itemWhiteSpace itemKind = iota
	itemEOF
	itemText
	itemError
)

const (
	whitespaceMarker = " "
)

type item struct {
	position int

	// kind signals how we've classified the data we have accumulated while
	// scanning the string.
	kind itemKind

	// value is the segment of data we've accumulated.
	value string
}

// Parse receives text and returns a string where the characters in a word have
// ben reversed; however, any whitespace was preserved (e.g., you can not simply
// reverse the entire string provided).
func Parse(text string) (string, error) {
	return text, nil
}

// stateFn is a function that is specific to a state within the string.
type stateFn func(*lexer) stateFn

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		state: lexText,
		items: make(chan item, 1),
	}

	go l.run()
}

type lexer struct {
	input    string    // we'll store the string being parsed
	start    int       // the position we started scanning
	position int       // the current position of our scan
	width    int       // we'll be using runes which can be double byte
	state    stateFn   // the current state function
	items    chan item // the channel we'll use to communicate between the lexer and the parser
}

func (l *lexer) emit(k itemKind) {
	i := item{
		position: l.start,
		kind:     k,
		value:    l.input[l.start:l.position],
	}

	l.items <- i

	l.ignore() // reset our scanner now that we've dispatched a segment
}

// ignore resets the start position to the current scan position effectively
// ignoring any input.
func (l *lexer) ignore() {
	l.start = l.position
}

func (l *lexer) next() (r rune) {
	if l.position >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.position:])
	l.position += l.width
	return r
}

func (l *lexer) backup() {
	l.position = l.position - 1
}

// run will scan the provided text and execute state functions.
func (l *lexer) run() {
	// We'll start, by assuming we're processing text.
	for fn := lexText; fn != nil; {
		fn = fn(l)
	}
	close(l.items)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	msg := fmt.Sprintf(format, args...)
	l.items <- item{
		kind:  itemError,
		value: msg,
	}

	return nil // nil stops the run
}

func lexText(l *lexer) stateFn {
	for {
		r := l.next()

		switch {
		case r == eof:
			l.backup()
			l.emit(itemText)
			l.emit(itemEOF)
			return nil
		case unicode.IsSpace(r):
			// emit any text we've accumulated.
			if l.position > l.start {
				l.emit(itemText)
			}

			return lexWhitespace
		}
	}
}

// lexWhitespace is the stateFn to run when the lexer has determined that we are
// now processing whitespace.
func lexWhitespace(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			l.backup()
			l.emit(itemWhiteSpace)
			l.emit(itemEOF)
			return nil
		case !unicode.IsSpace(r):
			if l.position > l.start {
				l.emit(itemWhiteSpace)
			}
			return lexText
		}
	}
}
