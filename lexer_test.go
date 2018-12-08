package parserlexer

import (
	"fmt"
	"testing"
)

type testcase struct {
	input  string
	expect string
}

func table() []testcase {
	return []testcase{
		{
			input:  "  which   is  harder to   do for     some examples    than others",
			expect: "  hcihw   si  redrah ot   od rof     emos selpmaxe    naht srehto",
		},
		{
			input:  "when we have a lot of words we want to preserve the white space",
			expect: "nehw ew evah a tol fo sdrow ew tnaw ot evreserp eht etihw ecaps",
		},
	}
}

// BenchmarkSimpleParser benchmarks the simple parser function which resulted in
// the following results on a 2014 MacBook Pro:
//
// BenchmarkSimpleParser-8          1000000              1119 ns/op             384 B/op         22 allocs/op
func BenchmarkSimpleParser(b *testing.B) {
	b.ReportAllocs()

	tc := table()[0]
	for n := 0; n < b.N; n++ {
		_ = SimpleParser(tc.input)
	}
}

func TestSimpleParser(t *testing.T) {
	tests := table()

	for i, test := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actual, err := Parse(test.input)
			if err != nil {
				t.Fatalf("no errors expected: %s", err)
			}
			if actual != test.expect {
				t.Fatalf("expected %q\nactual: %q\n", test.expect, actual)
			}
		})
	}
}

// BenchmarkParse benchmarks the parse function which resulted in the following
// results on a 2014 MacBook Pro:
//
// BenchmarkParse-8          200000              6360 ns/op             561 B/op         19 allocs/op
func BenchmarkLexParse(b *testing.B) {
	b.ReportAllocs()

	tc := table()[0]
	for n := 0; n < b.N; n++ {
		_, _ = Parse(tc.input)
	}
}

func TestLexParse(t *testing.T) {
	tests := table()

	for i, test := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actual, err := Parse(test.input)
			if err != nil {
				t.Fatalf("no errors expected: %s", err)
			}
			if actual != test.expect {
				t.Fatalf("expected %q\nactual: %q\n", test.expect, actual)
			}
		})
	}
}
