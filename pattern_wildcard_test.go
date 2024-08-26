package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

type testTable struct {
	input          string
	shouldMatch    []string
	shouldNotMatch []string
}

func TestCompilePatternWildcard(t *testing.T) {
	test := []testTable{
		{input: "ha*o", shouldMatch: []string{"hao", "halo", "hallo", "haweltlo"}, shouldNotMatch: []string{"", "welt", "hoa", "haa"}},
		{input: "wel*", shouldMatch: []string{"wel", "welt", "weltttttttt"}, shouldNotMatch: []string{"", "hallo", "walt", "wet"}},
		{input: "*elt", shouldMatch: []string{"welt", "elt", "weltttttttelt"}, shouldNotMatch: []string{"", "wel", "walt", "wet"}},
		{input: "*", shouldMatch: []string{"anything", ""}, shouldNotMatch: []string{}},
		{input: "*-mon-*", shouldMatch: []string{"s1-mon-test", "s1-mon-mon-mon-test"}, shouldNotMatch: []string{"se1-monn-test"}},
	}

	for _, tt := range test {
		start := newNfaStep()
		fm := compilePatternWildcard(start, charReplace(str2value(tt.input, nil, nil), charWildcard, byteWildcard), nil)

		for _, m := range tt.shouldMatch {
			target := transitionNfa(start, str2value(m, nil, nil), nil)
			assert.Check(t, len(target) > 0, "expected match '%s' with pattern '%s'", m, tt.input)
			assert.Check(t, fm == target[0], "expected match '%s' with pattern '%s'", m, tt.input)
		}

		for _, n := range tt.shouldNotMatch {
			target := transitionNfa(start, str2value(n, nil, nil), nil)
			assert.Check(t, len(target) == 0, "expected not to match '%s' with pattern '%s'", n, tt.input)
		}
	}
}

func TestValidatePatternWildcard_EmptyValue(t *testing.T) {
	pattern := &Pattern{
		Type: PatternWildcard,
	}
	err := validatePatternWildcard(pattern)
	assert.ErrorContains(t, err, "[wildcard] must contain a value")
}

func TestValidatePatternWildcard_SubPatterns(t *testing.T) {
	pattern := &Pattern{
		Type: PatternWildcard,
		Sub: []Pattern{
			{Type: PatternEquals, Value: "test"},
		},
	}
	err := validatePatternWildcard(pattern)
	assert.ErrorContains(t, err, "[wildcard] must contain a value")
}

func TestValidatePatternWildcard_SubPatterns2(t *testing.T) {
	pattern := &Pattern{
		Type:  PatternWildcard,
		Value: "invalid",
		Sub: []Pattern{
			{Type: PatternEquals, Value: "test"},
		},
	}
	err := validatePatternWildcard(pattern)
	assert.ErrorContains(t, err, "[wildcard] must not contain sub-patterns")
}

func TestValidatePatternWildcard_TwoConsecutiveWildcards(t *testing.T) {
	pattern := &Pattern{
		Type:  PatternWildcard,
		Value: "**test",
	}
	err := validatePatternWildcard(pattern)
	assert.ErrorContains(t, err, "[wildcard] must not contain two consecutive wildcards")
}
