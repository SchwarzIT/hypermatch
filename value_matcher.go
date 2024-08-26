package hypermatch

type ValueMatcher struct {
	Nfa *NfaStep `json:"s,omitempty"`
	// Maybe shortcuts here in the future
}

func newValueMatcher() *ValueMatcher {
	return &ValueMatcher{Nfa: newNfaStep()}
}

func (v *ValueMatcher) Transition(value []byte) []*FieldMatcher {
	return transitionNfa(v.Nfa, value, nil)
}
