// Package
package parserlexer

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "when we have a lot of words we want to preserve the white space",
			expect: "nehw ew evah a tol fo sdrow ew tnaw ot evreserp eht etihw ecaps",
		},
		{
			input:  "  which   is  harder to   do for     some examples    than others",
			expect: "  hcihw   si  redrah ot   od rof     emos selpmaxe    naht srehto",
		},
	}

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
