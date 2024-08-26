package hypermatch

import (
	"log"
	"strings"
)

func compileCondition(fm *FieldMatcher, id RuleIdentifier, cond *Condition) *FieldMatcher {
	if cond == nil || len(cond.Path) == 0 {
		return fm
	}

	return compilePattern(id, cond.Path, &cond.Pattern, fm, nil)
}

func str2value(input string, prefix []byte, suffix []byte) []byte {
	i := []byte(strings.ToLower(input))
	r := make([]byte, 0, len(prefix)+len(i)+len(suffix)+1)
	r = append(r, prefix...)
	r = append(r, i...)
	r = append(r, suffix...)
	r = append(r, byteValueTerminator)
	return r
}

func charReplace(data []byte, search byte, replace byte) []byte {
	for i := 0; i < len(data); i++ {
		if data[i] == search {
			data[i] = replace
		}
	}
	return data
}

func compilePattern(id RuleIdentifier, path string, pattern *Pattern, sourceFm *FieldMatcher, exitFm *FieldMatcher) *FieldMatcher {
	if pattern == nil || len(path) == 0 {
		return sourceFm
	}

	switch pattern.Type {
	case PatternEquals:
		return compilePatternEquals(sourceFm.GetTransition(path).Nfa, str2value(pattern.Value, nil, nil), exitFm)
	case PatternPrefix:
		return compilePatternWildcard(sourceFm.GetTransition(path).Nfa, str2value(pattern.Value, nil, []byte{byteWildcard}), exitFm)
	case PatternSuffix:
		return compilePatternWildcard(sourceFm.GetTransition(path).Nfa, str2value(pattern.Value, []byte{byteWildcard}, nil), exitFm)
	case PatternWildcard:
		return compilePatternWildcard(sourceFm.GetTransition(path).Nfa, charReplace(str2value(pattern.Value, nil, nil), charWildcard, byteWildcard), exitFm)
	case PatternAnythingBut:
		return compilePatternAnythingBut(id, path, pattern, sourceFm, exitFm)
	case PatternAnyOf:
		return compilePatternAnyOf(id, path, pattern, sourceFm, exitFm)
	case PatternAllOf:
		return compilePatternAllOf(id, path, pattern, sourceFm, exitFm)
	default:
		log.Printf("HyperMatch: unknown pattern type '%d', skipping!", pattern.Type)
	}

	return sourceFm
}
