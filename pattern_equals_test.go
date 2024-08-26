package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestCompilePatternEquals(t *testing.T) {
	test := []testTable{
		{input: "hallo", shouldMatch: []string{"hallo"}, shouldNotMatch: []string{"", "halloo", "halo"}},
		{input: "welt", shouldMatch: []string{"welt"}, shouldNotMatch: []string{"", "weltt"}},
	}

	for _, tt := range test {
		start := newNfaStep()
		fm := compilePatternEquals(start, str2value(tt.input, nil, nil), nil)

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
func TestValidatePatternEquals_EmptyValue(t *testing.T) {
	pattern := &Pattern{
		Type: PatternEquals,
	}

	err := validatePatternEquals(pattern)
	assert.ErrorContains(t, err, "[equals] must contain a value")
}

func TestValidatePatternEquals_SubPatterns(t *testing.T) {
	pattern := &Pattern{
		Type: PatternEquals,
		Sub: []Pattern{
			{Type: PatternEquals, Value: "subpattern"},
		},
	}

	err := validatePatternEquals(pattern)
	assert.ErrorContains(t, err, "[equals] must contain a value")
}
