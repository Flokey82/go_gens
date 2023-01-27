package genlanguage

import "testing"

func TestNumberToWords(t *testing.T) {
	for _, tc := range []struct {
		in  int
		out string
	}{
		{0, "zero"},
		{1, "one"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
		{5, "five"},
		{6, "six"},
		{7, "seven"},
		{8, "eight"},
		{9, "nine"},
		{10, "ten"},
		{11, "eleven"},
		{12, "twelve"},
		{13, "thirteen"},
		{14, "fourteen"},
		{15, "fifteen"},
		{16, "sixteen"},
		{17, "seventeen"},
		{18, "eighteen"},
		{19, "nineteen"},
		{20, "twenty"},
		{21, "twenty-one"},
		{22, "twenty-two"},
		{23, "twenty-three"},
		{24, "twenty-four"},
		{100, "one hundred"},
		{101, "one hundred and one"},
		{10933, "ten thousand nine hundred and thirty-three"},
		{1000000, "one million"},
		{1000001, "one million and one"},
	} {
		if got := NumberToWords(tc.in); got != tc.out {
			t.Errorf("NumberToWords(%d) = %q, want %q", tc.in, got, tc.out)
		}
	}
}
