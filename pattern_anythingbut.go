package hypermatch

import "fmt"

func validatePatternAnythingBut(path string, pattern *Pattern) error {
	if len(pattern.Value) == 0 && len(pattern.Sub) == 0 {
		return fmt.Errorf("[%s] must contain a value or sub-patterns", PatternAnythingBut.String())
	}

	for _, p := range pattern.Sub {
		if err := validatePattern(path, &p); err != nil {
			return err
		}
	}

	return nil
}

func compilePatternAnythingBut(id RuleIdentifier, path string, pattern *Pattern, sourceFm *fieldMatcher, exitFm *fieldMatcher) *fieldMatcher {
	exitFm = compilePatternAnyOf(id, path, pattern, sourceFm, exitFm)
	exitFm.MatchingAnythingButRuleIdentifiers = append(exitFm.MatchingAnythingButRuleIdentifiers, id)

	fm := newFieldMatcher()
	sourceFm.AddAnythingButTransition(id, path, fm)
	return fm
}
