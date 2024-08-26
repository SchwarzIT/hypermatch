package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

type compoundTestTable struct {
	pattern        Pattern
	shouldMatch    []string
	shouldNotMatch []string
}

func TestCompilePatternAnythingBut(t *testing.T) {
	test := []compoundTestTable{
		{pattern: Pattern{Type: PatternAnythingBut, Sub: []Pattern{{Type: PatternEquals, Value: "hallo"}}}, shouldMatch: []string{"welt"}, shouldNotMatch: []string{"hallo"}},
		{pattern: Pattern{Type: PatternAnythingBut, Sub: []Pattern{{Type: PatternWildcard, Value: "ha*o"}}}, shouldMatch: []string{"welt"}, shouldNotMatch: []string{"hallo"}},
	}

	for i, tt := range test {
		sourceFm := newFieldMatcher()
		compilePatternAnythingBut(RuleIdentifier(i), "test", &tt.pattern, sourceFm, nil)

		for _, m := range tt.shouldMatch {
			target := transitionNfa(sourceFm.GetTransition("test").Nfa, str2value(m, nil, nil), nil)
			assert.Check(t, len(target) == 0, "expected match '%s' with pattern '%v'", m, tt.pattern)
		}

		for _, n := range tt.shouldNotMatch {
			target := transitionNfa(sourceFm.GetTransition("test").Nfa, str2value(n, nil, nil), nil)
			assert.Check(t, len(target) > 0, "expected not to match '%s' with pattern '%v'", n, tt.pattern)
			_, ok := sourceFm.AnythingButTransitions["test"][target[0].MatchingAnythingButRuleIdentifiers[0]]
			assert.Check(t, ok, "expected not to match '%s' with pattern '%v'", n, tt.pattern)
		}
	}
}

func TestValidatePatternAnythingBut_EmptyValueAndSubPatterns(t *testing.T) {
	path := "testPath"
	pattern := &Pattern{Type: PatternAnythingBut}

	err := validatePatternAnythingBut(path, pattern)
	assert.ErrorContains(t, err, "[anythingBut] must contain a value or sub-patterns")
}

func TestValidatePatternAnythingBut_WithValueAndEmptySubPatterns(t *testing.T) {
	path := "testPath"
	pattern := &Pattern{Type: PatternAnythingBut, Value: "testValue"}

	err := validatePatternAnythingBut(path, pattern)
	assert.NilError(t, err)
}

func TestValidatePatternAnythingBut_EmptyValueAndNonEmptySubPatterns(t *testing.T) {
	path := "testPath"
	pattern := &Pattern{Type: PatternAnythingBut, Sub: []Pattern{
		{Type: PatternEquals, Value: "testValue"},
	}}

	err := validatePatternAnythingBut(path, pattern)
	assert.NilError(t, err)
}
