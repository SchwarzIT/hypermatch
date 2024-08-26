package hypermatch

import (
	"gotest.tools/v3/assert"
	"log"
	"slices"
	"testing"
)

func containsAll(data []RuleIdentifier, entries ...RuleIdentifier) bool {
	for _, e := range entries {
		if !slices.Contains(data, e) {
			log.Printf("Does not contain %d: %v", e, data)
			return false
		}
	}

	return true
}

func TestHyperMatchSimple(t *testing.T) {

	h := NewHyperMatch()

	assert.NilError(t, h.AddRule(1, []Condition{
		{Path: "y", Pattern: Pattern{Type: PatternEquals, Value: "bb"}},
		{Path: "x", Pattern: Pattern{Type: PatternEquals, Value: "aa"}},
	}))

	assert.NilError(t, h.AddRule(2, []Condition{
		{Path: "y", Pattern: Pattern{Type: PatternEquals, Value: "bb"}},
	}))

	assert.NilError(t, h.AddRule(3, []Condition{
		{Path: "x", Pattern: Pattern{Type: PatternEquals, Value: "dd"}},
		{Path: "y", Pattern: Pattern{Type: PatternEquals, Value: "cc"}},
	}))

	assert.NilError(t, h.AddRule(4, []Condition{
		{Path: "x", Pattern: Pattern{Type: PatternPrefix, Value: "a"}},
		{Path: "y", Pattern: Pattern{Type: PatternPrefix, Value: "b"}},
	}))

	assert.NilError(t, h.AddRule(5, []Condition{
		{Path: "x", Pattern: Pattern{Type: PatternSuffix, Value: "x"}},
		{Path: "y", Pattern: Pattern{Type: PatternSuffix, Value: "y"}},
	}))

	assert.NilError(t, h.AddRule(6, []Condition{
		{Path: "y", Pattern: Pattern{Type: PatternAnythingBut, Sub: []Pattern{
			{Type: PatternEquals, Value: "bb"},
		}}},
	}))

	assert.NilError(t, h.AddRule(7, []Condition{
		{Path: "y", Pattern: Pattern{Type: PatternAnythingBut, Sub: []Pattern{
			{Type: PatternEquals, Value: "b*y"},
		}}},
	}))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "x", Values: []string{"aa"}},
		{Path: "y", Values: []string{"bb"}},
	}), 1, 2, 4))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "x", Values: []string{"dd"}},
		{Path: "y", Values: []string{"cc"}},
	}), 3, 6))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "x", Values: []string{"aax"}},
		{Path: "y", Values: []string{"bby"}},
	}), 4, 5, 6, 7))

	assert.Check(t, containsAll(h.Match([]Property{
		{Path: "x", Values: []string{"aax", "dd"}},
		{Path: "y", Values: []string{"bby", "cc"}},
	}), 3, 4, 5, 6, 7))

}
func TestValidateRule(t *testing.T) {
	tests := []struct {
		name    string
		set     ConditionSet
		wantErr bool
	}{
		{
			name:    "empty set",
			set:     []Condition{},
			wantErr: false,
		},
		{
			name: "valid set",
			set: []Condition{
				{
					Path: "a",
					Pattern: Pattern{
						Type:  PatternEquals,
						Value: "2",
					},
				},
				{
					Path: "b",
					Pattern: Pattern{
						Type:  PatternEquals,
						Value: "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid set",
			set: []Condition{
				{
					Path: "a",
					Pattern: Pattern{
						Type: PatternEquals,
						Sub: []Pattern{
							{
								Type:  PatternEquals,
								Value: "2",
							},
						},
					},
				},
				{
					Path: "b",
					Pattern: Pattern{
						Type:  PatternEquals,
						Value: "1",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRule(tt.set); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
