package hypermatch

type FieldMatcher struct {
	Transitions                        map[string]*ValueMatcher                    `json:"t,omitempty"`
	MatchingRuleIdentifiers            []RuleIdentifier                            `json:"m,omitempty"`
	MatchingAnythingButRuleIdentifiers []RuleIdentifier                            `json:"n,omitempty"`
	AnythingButTransitions             map[string]map[RuleIdentifier]*FieldMatcher `json:"o,omitempty"`
	Exclusive                          bool                                        `json:"e,omitempty"`
}

func newFieldMatcher() *FieldMatcher {
	return &FieldMatcher{
		Transitions:                        make(map[string]*ValueMatcher),
		MatchingRuleIdentifiers:            nil,
		MatchingAnythingButRuleIdentifiers: nil,
		AnythingButTransitions:             nil,
	}
}

func (f *FieldMatcher) AddAnythingButTransition(id RuleIdentifier, path string, fm *FieldMatcher) {
	if f.AnythingButTransitions == nil {
		f.AnythingButTransitions = make(map[string]map[RuleIdentifier]*FieldMatcher)
	}
	if _, ok := f.AnythingButTransitions[path]; !ok {
		f.AnythingButTransitions[path] = make(map[RuleIdentifier]*FieldMatcher)
	}
	f.AnythingButTransitions[path][id] = fm
}

func (f *FieldMatcher) GetTransition(key string) *ValueMatcher {
	vm, ok := f.Transitions[key]
	if !ok {
		vm = newValueMatcher()
		f.Transitions[key] = vm
	}
	return vm
}
