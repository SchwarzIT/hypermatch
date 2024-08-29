package hypermatch

type valueMatcher struct {
	Nfa *nfaStep `json:"s,omitempty"`
	// Maybe shortcuts here in the future
}

func newValueMatcher() *valueMatcher {
	return &valueMatcher{Nfa: newNfaStep()}
}

func (v *valueMatcher) Transition(value []byte) []*fieldMatcher {
	return transitionNfa(v.Nfa, value, nil)
}
