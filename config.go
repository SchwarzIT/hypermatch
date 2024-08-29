package hypermatch

// RuleIdentifier is a type alias to represent the identifier for a rule, so that the user of hypermatch get identify which of the rules matches an event.
type RuleIdentifier any

const (
	byteValueTerminator byte = 0xf5
	byteWildcard        byte = 0xf6

	charWildcard byte = '*'
)
