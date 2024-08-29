package hypermatch

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type HyperMatch struct {
	matcher    *fieldMatcher
	rulesCount uint64
}

func NewHyperMatch() *HyperMatch {
	return &HyperMatch{matcher: newFieldMatcher(), rulesCount: 0}
}

// ValidateRule validates the given condition set.
// Returns an error if validation fails.
func ValidateRule(set ConditionSet) error {
	for _, s := range set {
		if err := validateCondition(&s); err != nil {
			return errors.Join(fmt.Errorf("could not validate condition for '%s'", s.Path), err)
		}
	}
	return nil
}

// AddRule adds a new rule to the HyperMatch instance.
//
// The rule is defined by a unique identifier and a set of conditions.
// The conditions define the properties that must be matched in order for the rule to be triggered.
//
// If the rule is successfully added, the function returns nil.
// Otherwise, an error is returned.
func (m *HyperMatch) AddRule(id RuleIdentifier, conditionSet ConditionSet) error {
	if len(conditionSet) == 0 {
		return fmt.Errorf("no conditions provided")
	}

	sort.Slice(conditionSet, func(i, j int) bool {
		return strings.ToLower(conditionSet[i].Path) < strings.ToLower(conditionSet[j].Path)
	})

	cfm := m.matcher
	for _, s := range conditionSet {
		cfm = compileCondition(cfm, id, &s)
	}

	cfm.MatchingRuleIdentifiers = append(cfm.MatchingRuleIdentifiers, id)

	m.rulesCount += 1

	return nil
}

// Match takes a list of properties and returns a list of rule identifiers that match those properties
func (m *HyperMatch) Match(properties []Property) []RuleIdentifier {

	sort.Slice(properties, func(i, j int) bool {
		return strings.ToLower(properties[i].Path) < strings.ToLower(properties[j].Path)
	})

	matches := newMatchSet()

	for i := range properties {
		tryToMatch(properties, i, m.matcher, matches)
	}

	return matches.All()
}

func tryToMatch(properties []Property, i int, fm *fieldMatcher, set *matchSet) {
	if i >= len(properties) {
		return
	}
	field := properties[i]

	nextMatchers := match(fm, field.Path, field.Values)
	for _, m := range nextMatchers {
		set = set.Add(m.MatchingRuleIdentifiers...)
		for nI := i; nI < len(properties); nI++ {
			tryToMatch(properties, nI, m, set)
		}
	}
}

func match(f *fieldMatcher, field string, values []string) []*fieldMatcher {
	vm, ok := f.Transitions[field]
	if !ok {
		return nil
	}

	var afms []*fieldMatcher

	for _, value := range values {
		v := str2value(value, nil, nil)
		fms := vm.Transition(v)

		/*	if len(fms) == 0 {
				continue
			}
		*/
		set := newMatchSet()
		for _, f := range fms {
			set.Add(f.MatchingAnythingButRuleIdentifiers...)
		}

		if ts, ok := f.AnythingButTransitions[field]; ok {
			for id, fm := range ts {
				if !set.Contains(id) {
					fms = append(fms, fm)
				}
			}
		}

		afms = append(afms, fms...)
	}

	return afms
}
