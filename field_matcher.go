package hypermatch

type fieldMatcher struct {
	Transitions                        map[string]*valueMatcher                    `json:"t,omitempty"`
	MatchingRuleIdentifiers            []RuleIdentifier                            `json:"m,omitempty"`
	MatchingAnythingButRuleIdentifiers []RuleIdentifier                            `json:"n,omitempty"`
	AnythingButTransitions             map[string]map[RuleIdentifier]*fieldMatcher `json:"o,omitempty"`
	Exclusive                          bool                                        `json:"e,omitempty"`
}

func newFieldMatcher() *fieldMatcher {
	return &fieldMatcher{
		Transitions:                        make(map[string]*valueMatcher),
		MatchingRuleIdentifiers:            nil,
		MatchingAnythingButRuleIdentifiers: nil,
		AnythingButTransitions:             nil,
	}
}

func (f *fieldMatcher) AddAnythingButTransition(id RuleIdentifier, path string, fm *fieldMatcher) {
	if f.AnythingButTransitions == nil {
		f.AnythingButTransitions = make(map[string]map[RuleIdentifier]*fieldMatcher)
	}
	if _, ok := f.AnythingButTransitions[path]; !ok {
		f.AnythingButTransitions[path] = make(map[RuleIdentifier]*fieldMatcher)
	}
	f.AnythingButTransitions[path][id] = fm
}

func (f *fieldMatcher) GetTransition(key string) *valueMatcher {
	vm, ok := f.Transitions[key]
	if !ok {
		vm = newValueMatcher()
		f.Transitions[key] = vm
	}
	return vm
}
