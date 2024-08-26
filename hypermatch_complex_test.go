package hypermatch

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestHyperMatchComplex(t *testing.T) {

	h := NewHyperMatch()

	assert.NilError(t, h.AddRule(1, []Condition{
		{Path: "namespace", Pattern: Pattern{Type: PatternEquals, Value: "monitoring-schwarz"}},
		{Path: "l1", Pattern: Pattern{Type: PatternEquals, Value: "lidl"}},
		{Path: "name", Pattern: Pattern{Type: PatternWildcard, Value: "*OS*"}},
		{Path: "priority", Pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{{Type: PatternEquals, Value: "P1"}, {Type: PatternEquals, Value: "P2"}}}},
		{Path: "hostgroups", Pattern: Pattern{Type: PatternAllOf, Sub: []Pattern{{Type: PatternEquals, Value: "lidl-de"}, {Type: PatternEquals, Value: "store-servers"}}}},
	}))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "namespace", Values: []string{"monitoring-schwarz"}},
		{Path: "l1", Values: []string{"lidl"}},
		{Path: "name", Values: []string{"OS hdd"}},
		{Path: "priority", Values: []string{"P1"}},
		{Path: "hostgroups", Values: []string{"lidl", "lidl-de", "de", "store-servers"}},
	},
	), 1))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "namespace", Values: []string{"monitoring-schwarz"}},
		{Path: "l1", Values: []string{"lidl"}},
		{Path: "name", Values: []string{"OS hdd"}},
		{Path: "priority", Values: []string{"P1"}},
		{Path: "hostgroups", Values: []string{"lidl", "lidl-de", "de"}},
	},
	)))
}
