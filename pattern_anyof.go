package hypermatch

import "fmt"

func validatePatternAnyOf(path string, pattern *Pattern) error {
	if len(pattern.Value) > 0 {
		return fmt.Errorf("[%s] must not contain a value", PatternAnyOf.String())
	}
	if len(pattern.Sub) == 0 {
		return fmt.Errorf("[%s] must contain sub-patterns", PatternAnyOf.String())
	}

	for _, p := range pattern.Sub {
		if err := validatePattern(path, &p); err != nil {
			return err
		}
	}

	return nil
}

func compilePatternAnyOf(id RuleIdentifier, path string, pattern *Pattern, sourceFm *fieldMatcher, exitFm *fieldMatcher) *fieldMatcher {
	if exitFm != nil && !exitFm.Exclusive {
		exitFm.Exclusive = true
	}
	for _, p := range pattern.Sub {
		exitFm = compilePattern(id, path, &p, sourceFm, exitFm)
		exitFm.Exclusive = true
	}

	return exitFm
}
