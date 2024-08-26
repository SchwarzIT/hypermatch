package hypermatch

type RuleIdentifier any

const (
	byteValueTerminator byte = 0xf5
	byteWildcard        byte = 0xf6

	charWildcard byte = '*'
)
