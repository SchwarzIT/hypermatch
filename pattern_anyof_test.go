package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestCompilePatternAnyOf(t *testing.T) {
	test := []compoundTestTable{
		{pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternEquals, Value: "welt"},
		}}, shouldMatch: []string{"welt", "hallo"}, shouldNotMatch: []string{"was"}},
		{pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternWildcard, Value: "wel*"},
		}}, shouldMatch: []string{"welt", "hallo", "weltttt"}, shouldNotMatch: []string{"was"}},
	}

	for i, tt := range test {
		sourceFm := newFieldMatcher()
		fm := compilePatternAnyOf(RuleIdentifier(i), "test", &tt.pattern, sourceFm, nil)

		for _, m := range tt.shouldMatch {
			target := transitionNfa(sourceFm.GetTransition("test").Nfa, str2value(m, nil, nil), nil)
			assert.Check(t, len(target) > 0, "expected match '%s' with pattern '%v'", m, tt.pattern)
			assert.Check(t, fm == target[0], "expected match '%s' with pattern '%v'", m, tt.pattern)
		}

		for _, n := range tt.shouldNotMatch {
			target := transitionNfa(sourceFm.GetTransition("test").Nfa, str2value(n, nil, nil), nil)
			assert.Check(t, len(target) == 0, "expected not to match '%s' with pattern '%v'", n, tt.pattern)
		}
	}
}

func TestValidatePatternAnyOf_ErrorForValue(t *testing.T) {
	pattern := Pattern{
		Type:  PatternAnyOf,
		Value: "invalid",
		Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternEquals, Value: "welt"},
		},
	}

	err := validatePatternAnyOf("test", &pattern)
	assert.ErrorContains(t, err, "[anyOf] must not contain a value")
}
func TestValidatePatternAnyOf_ErrorForNoSubPatterns(t *testing.T) {
	pattern := Pattern{
		Type: PatternAnyOf,
		Sub:  []Pattern{},
	}

	err := validatePatternAnyOf("test", &pattern)
	assert.ErrorContains(t, err, "[anyOf] must contain sub-patterns")
}
func TestValidatePatternAnyOf_CorrectlyValidatesSubPatterns(t *testing.T) {
	pattern := Pattern{
		Type: PatternAnyOf,
		Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternEquals, Value: "welt"},
		},
	}

	err := validatePatternAnyOf("test", &pattern)
	assert.NilError(t, err)
}
