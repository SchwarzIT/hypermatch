package hypermatch

import "fmt"

func validatePatternAllOf(path string, pattern *Pattern) error {
	if len(pattern.Value) > 0 {
		return fmt.Errorf("'[%s] must not contain a value", PatternAllOf.String())
	}
	if len(pattern.Sub) == 0 {
		return fmt.Errorf("[%s] must contain sub-patterns", PatternAllOf.String())
	}

	for _, p := range pattern.Sub {
		if err := validatePattern(path, &p); err != nil {
			return err
		}
	}

	return nil
}

func compilePatternAllOf(id RuleIdentifier, path string, pattern *Pattern, sourceFm *fieldMatcher, exitFm *fieldMatcher) *fieldMatcher {
	lastSourceFm := sourceFm

	for i, p := range pattern.Sub {
		if i == len(pattern.Sub)-1 {
			// the last pattern
			exitFm = compilePattern(id, path, &p, lastSourceFm, exitFm)
		} else {
			lastSourceFm = compilePattern(id, path, &p, lastSourceFm, nil)
		}
	}

	return exitFm
}
