package hypermatch

type matchSet struct {
	// this is faster and more memory efficient than map[RuleIdentifier]bool !
	set map[RuleIdentifier]struct{}
}

func newMatchSet() *matchSet {
	return &matchSet{set: make(map[RuleIdentifier]struct{})}
}

func (m *matchSet) Add(matches ...RuleIdentifier) *matchSet {
	for _, x := range matches {
		m.set[x] = struct{}{}
	}
	return m
}

func (m *matchSet) Contains(r RuleIdentifier) bool {
	_, ok := m.set[r]
	return ok
}

func (m *matchSet) All() []RuleIdentifier {
	matches := make([]RuleIdentifier, 0, len(m.set))
	for x := range m.set {
		matches = append(matches, x)
	}
	return matches
}
