package hypermatch

import (
	"fmt"
)

func validateCondition(cond *Condition) error {
	if cond == nil || len(cond.Path) == 0 {
		return nil
	}

	return validatePattern(cond.Path, &cond.Pattern)
}

func validatePattern(path string, pattern *Pattern) error {
	if pattern == nil || len(path) == 0 {
		return nil
	}

	switch pattern.Type {
	case PatternEquals:
		return validatePatternEquals(pattern)
	case PatternPrefix:
		return validatePatternWildcard(pattern)
	case PatternSuffix:
		return validatePatternWildcard(pattern)
	case PatternWildcard:
		return validatePatternWildcard(pattern)
	case PatternAnythingBut:
		return validatePatternAnythingBut(path, pattern)
	case PatternAnyOf:
		return validatePatternAnyOf(path, pattern)
	case PatternAllOf:
		return validatePatternAllOf(path, pattern)
	default:
		return fmt.Errorf("unknown pattern type '%d'", pattern.Type)
	}

}
