package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

type compoundMultiTestTable struct {
	pattern        Pattern
	shouldMatch    [][]string
	shouldNotMatch [][]string
}

func TestCompilePatternAllOf(t *testing.T) {
	test := []compoundMultiTestTable{
		{pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternEquals, Value: "welt"},
		}}, shouldMatch: [][]string{{"welt", "hallo"}}, shouldNotMatch: [][]string{{"welt"}, {"hallo"}, {"welt", "eins"}}},
		{pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{
			{Type: PatternEquals, Value: "hallo"},
			{Type: PatternWildcard, Value: "wel*"},
		}}, shouldMatch: [][]string{{"welt", "hallo", "weltttt"}}, shouldNotMatch: [][]string{{"welt"}, {"hallo"}, {"welt", "eins"}}},
	}

	for i, tt := range test {
		sourceFm := newFieldMatcher()
		fm := compilePatternAllOf(RuleIdentifier(i), "test", &tt.pattern, sourceFm, nil)

		for _, m := range tt.shouldMatch {
			start := sourceFm
			for range m {
				match := matchAny(start, m)
				assert.Check(t, match != nil, "expected match '%s' with pattern '%v'", m, tt.pattern)
				start = match
				if start == fm {
					break
				}
			}
			assert.Check(t, fm == start, "expected match '%s' with pattern '%v'", m, tt.pattern)
		}

		for _, n := range tt.shouldNotMatch {
			start := sourceFm
			found := true
			for range n {
				match := matchAny(start, n)
				if match == nil {
					found = false
					break
				}
				start = match
			}
			assert.Check(t, !found || fm != start, "expected not to match '%s' with pattern '%v'", n, tt.pattern)
		}
	}
}

func matchAny(start *fieldMatcher, values []string) *fieldMatcher {
	for _, v := range values {
		fm := transitionNfa(start.GetTransition("test").Nfa, str2value(v, nil, nil), nil)
		if len(fm) > 0 {
			return fm[0]
		}
	}
	return nil
}

func TestValidatePatternAllOf_WithValue(t *testing.T) {
	err := validatePatternAllOf("test", &Pattern{
		Type:  PatternAllOf,
		Value: "invalid",
		Sub: []Pattern{
			{Type: PatternEquals, Value: "value1"},
			{Type: PatternEquals, Value: "value2"},
		},
	})
	assert.ErrorContains(t, err, "'[allOf] must not contain a value")
}
func TestValidatePatternAllOf_NoSubPatterns(t *testing.T) {
	err := validatePatternAllOf("test", &Pattern{
		Type: PatternAllOf,
		Sub:  []Pattern{},
	})
	assert.ErrorContains(t, err, "[allOf] must contain sub-patterns")
}

func TestValidatePatternAllOf_NestedPatterns(t *testing.T) {
	err := validatePatternAllOf("test", &Pattern{
		Type: PatternAllOf,
		Sub: []Pattern{
			{Type: PatternEquals, Value: "value1"},
			{
				Type: PatternAllOf,
				Sub: []Pattern{
					{Type: PatternEquals, Value: "value2"},
					{Type: PatternEquals, Value: "value3"},
				},
			},
			{Type: PatternEquals, Value: "value4"},
		},
	})
	assert.NilError(t, err)
}
