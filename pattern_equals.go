package hypermatch

import "fmt"

func validatePatternEquals(pattern *Pattern) error {
	if len(pattern.Value) == 0 {
		return fmt.Errorf("[%s] must contain a value", PatternEquals.String())
	}
	if len(pattern.Sub) > 0 {
		return fmt.Errorf("[%s] must not contain sub-patterns", PatternEquals.String())
	}
	return nil
}

func compilePatternEquals(start *nfaStep, value []byte, exitFm *fieldMatcher) *fieldMatcher {
	step := start
	for _, char := range value {
		step = step.MakeStep(char)
	}

	return step.addOrReuseOrCreateFieldTransition(exitFm)
}
