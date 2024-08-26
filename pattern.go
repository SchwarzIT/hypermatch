package hypermatch

import (
	"strings"
)

type PatternType int

const (
	PatternEquals PatternType = iota
	PatternPrefix
	PatternSuffix
	PatternWildcard

	PatternAnythingBut
	PatternAnyOf
	PatternAllOf

	PatternUnknown
)

func (p PatternType) AllValues() []PatternType {
	return []PatternType{PatternEquals, PatternPrefix, PatternSuffix, PatternWildcard, PatternAnythingBut, PatternAnyOf, PatternAllOf}
}

func (p PatternType) HasLiteralValue() bool {
	switch p {
	case PatternEquals, PatternPrefix, PatternSuffix, PatternWildcard:
		return true
	default:
		return false
	}
}

func (p PatternType) String() string {
	switch p {
	case PatternEquals:
		return "equals"
	case PatternPrefix:
		return "prefix"
	case PatternSuffix:
		return "suffix"
	case PatternWildcard:
		return "wildcard"
	case PatternAnythingBut:
		return "anythingBut"
	case PatternAnyOf:
		return "anyOf"
	case PatternAllOf:
		return "allOf"
	default:
		return ""
	}
}

func PatternTypeFromString(input string) PatternType {
	switch strings.ToLower(input) {
	case "equals":
		return PatternEquals
	case "prefix":
		return PatternPrefix
	case "suffix":
		return PatternSuffix
	case "wildcard":
		return PatternWildcard
	case "anythingbut":
		return PatternAnythingBut
	case "anyof":
		return PatternAnyOf
	case "allof":
		return PatternAllOf
	}
	return PatternUnknown
}
