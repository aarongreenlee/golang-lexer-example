package parserlexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Parse(text string) (string, error) {
	fmt.Printf("partsing: %q\n", text)

	p := parser{
		lex: lex(text),
	}

	p.parse()

	if p.errItem != nil {
		return "", fmt.Errorf("error processing the following %q", p.errItem.value)
	}

	return p.result, nil
}

type parser struct {
	result  string
	lex     *lexer
	errItem *item
}

func (p *parser) parse() {
	sb := strings.Builder{}

	for item := range p.lex.items {
		switch item.kind {
		case itemEOF:
			p.result = sb.String()
			return
		case itemError:
			p.errItem = &item
			return
		case itemText:
			sb.WriteString(reverse(item.value))

		case itemWhiteSpace:
			sb.WriteString(item.value)
		default:
			fmt.Printf("unknown type %q", item.kind)
		}
	}
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

type itemKind int

const (
	itemWhiteSpace itemKind = iota
	itemEOF
	itemText
	itemError
)

// item is accumulated while lexing the provided input, and emitted over a
// channel to the parser. Items could also be called tokens as we tokenize the
// input.
type item struct {
	position int

	// kind signals how we've classified the data we have accumulated while
	// scanning the string.
	kind itemKind

	// value is the segment of data we've accumulated.
	value string
}

const eof = -1

// stateFn is a function that is specific to a state within the string.
type stateFn func(*lexer) stateFn

// lex creates a lexer and starts scanning the provided input.
func lex(input string) *lexer {
	l := &lexer{
		input: input,
		state: lexText,
		items: make(chan item, 1),
	}

	go l.scan()

	return l
}

// lexer is created to manage an individual scanning/parsing operation.
type lexer struct {
	input    string    // we'll store the string being parsed
	start    int       // the position we started scanning
	position int       // the current position of our scan
	width    int       // we'll be using runes which can be double byte
	state    stateFn   // the current state function
	items    chan item // the channel we'll use to communicate between the lexer and the parser
}

// emit sends a item over the channel so the parser can collect and manage
// each segment.
func (l *lexer) emit(k itemKind) {

	// fmt.Printf("%d %d %d\n", l.start, l.position, len(l.input))

	accumulation := l.input[l.start:l.position]

	fmt.Printf("emitting %q\n", accumulation)

	i := item{
		position: l.start,
		kind:     k,
		value:    accumulation,
	}

	l.items <- i

	l.ignore() // reset our scanner now that we've dispatched a segment
}

// nextItem pulls an item from the lexer's result channel.
func (l *lexer) nextItem() item {
	return <-l.items
}

// ignore resets the start position to the current scan position effectively
// ignoring any input.
func (l *lexer) ignore() {
	l.start = l.position
}

// next advances the lexer state to the next rune.
func (l *lexer) next() (r rune) {
	if l.position >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.position:])
	l.position += l.width
	return r
}

// backup allows us to step back one run1e which is helpful when you've crossed
// a boundary from one state to another.
func (l *lexer) backup() {
	l.position = l.position - 1
}

// scan will step through the provided text and execute state functions as
// state changes are observed in the provided input.
func (l *lexer) scan() {
	// When we begin processing, let's assume we're going to process text.
	// One state function will return another until `nil` is returned to signal
	// the end of our process.
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

	return nil
}

// lexEOF emits the accumulated data classified by the provided itemKind and
// signals that we've reached the end of our lexing by returning `nil` instead
// of a state function.
func (l *lexer) lexEOF(k itemKind) stateFn {

	//	l.backup()
	if l.start > l.position {
		l.ignore()
	}

	l.emit(k)
	l.emit(itemEOF)
	return nil
}

// lexText scans what is expected to be text.
func lexText(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			return l.lexEOF(itemText)
		case unicode.IsSpace(r):
			l.backup()

			// emit any text we've accumulated.
			if l.position > l.start {
				l.emit(itemText)
			}
			return lexWhitespace
		}
	}
}

// lexWhitespace scans what is expected to be whitespace.
func lexWhitespace(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case r == eof:
			return l.lexEOF(itemWhiteSpace)
		case !unicode.IsSpace(r):
			l.backup()
			if l.position > l.start {
				l.emit(itemWhiteSpace)
			}
			return lexText
		}
	}
}
