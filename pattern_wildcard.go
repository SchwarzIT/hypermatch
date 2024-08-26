package hypermatch

import (
	"fmt"
	"strings"
)

func validatePatternWildcard(pattern *Pattern) error {
	if len(pattern.Value) == 0 {
		return fmt.Errorf("[%s] must contain a value", PatternWildcard.String())
	}
	if len(pattern.Sub) > 0 {
		return fmt.Errorf("[%s] must not contain sub-patterns", PatternWildcard.String())
	}

	if strings.Contains(pattern.Value, "**") {
		return fmt.Errorf("[%s] must not contain two consecutive wildcards", PatternWildcard.String())
	}

	return nil
}

func compilePatternWildcard(start *NfaStep, value []byte, exitFm *FieldMatcher) *FieldMatcher {
	step := start
	var lastWildcardStep *NfaStep
	for i, char := range value {
		if char == byteWildcard {
			if (i == len(value)-2 && value[len(value)-1] == byteValueTerminator) || i == len(value)-1 {
				// wildcard is the last character before value terminator
				return step.addOrReuseOrCreateFieldTransition(exitFm)
			} else {
				lastWildcardStep = step
				step = step.MakeStep(byteWildcard)
				step.ValueTransitions[byteWildcard] = step
			}
		} else {
			step = step.MakeStep(char)
			if lastWildcardStep != nil {
				lastWildcardStep.ValueTransitions[char] = step
				lastWildcardStep = nil
			}
		}
	}

	return step.addOrReuseOrCreateFieldTransition(exitFm)
}
