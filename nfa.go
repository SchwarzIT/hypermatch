package hypermatch

import "slices"

type NfaStep struct {
	ValueTransitions map[byte]*NfaStep `json:"v,omitempty"`
	FieldTransition  []*FieldMatcher   `json:"f,omitempty"`
}

func newNfaStep() *NfaStep {
	return &NfaStep{
		ValueTransitions: make(map[byte]*NfaStep),
		FieldTransition:  nil,
	}
}

func (n *NfaStep) MakeStep(char byte) *NfaStep {
	s, ok := n.ValueTransitions[char]
	if !ok {
		s = newNfaStep()
		n.ValueTransitions[char] = s
	}
	return s
}

func (n *NfaStep) addOrReuseOrCreateFieldTransition(fm *FieldMatcher) *FieldMatcher {
	if fm != nil {
		if slices.Index(n.FieldTransition, fm) == -1 {
			n.FieldTransition = append(n.FieldTransition, fm)
			return fm
		} else {
			return fm
		}
	} else {
		for _, f := range n.FieldTransition {
			if !f.Exclusive {
				return f
			}
		}
		fm := newFieldMatcher()
		n.FieldTransition = append(n.FieldTransition, fm)
		return fm
	}
}

func transitionNfa(step *NfaStep, value []byte, transitions []*FieldMatcher) []*FieldMatcher {
	if len(value) == 0 {
		return nil
	}

	if step.FieldTransition != nil {
		transitions = append(transitions, step.FieldTransition...)
	}

	// transition through the NFA
	for i, v := range value {
		// if there are no value transitions we can stop here
		if step.ValueTransitions == nil {
			break
		}

		// if there is a wildcard transition, run through it recursively
		if w, ok := step.ValueTransitions[byteWildcard]; ok && len(value) > i+1 {
			transitions = append(transitions, transitionNfa(w, value[i+1:], nil)...)
		}

		// follow the next step if possible, otherwise stop
		if s, ok := step.ValueTransitions[v]; ok {
			step = s
		} else {
			break
		}

		// if there are field transitions, add them!
		if step.FieldTransition != nil {
			transitions = append(transitions, step.FieldTransition...)
		}
	}

	return transitions
}
