package hypermatch

import (
	"encoding/json"
	"gotest.tools/v3/assert"
	"log"
	"testing"
)

func TestConditions_MarshalJSON(t *testing.T) {
	r := ConditionSet{
		{
			Path: "name",
			Pattern: Pattern{
				Type: PatternAllOf,
				Sub: []Pattern{
					{Type: PatternEquals, Value: "hallo"},
					{Type: PatternWildcard, Value: "hallo*"},
					{Type: PatternAnythingBut, Sub: []Pattern{
						{Type: PatternSuffix, Value: "test"},
						{Type: PatternPrefix, Value: "st"},
					}},
					{Type: PatternAnyOf, Sub: []Pattern{
						{Type: PatternSuffix, Value: "te"},
						{Type: PatternPrefix, Value: "tet"},
					}},
				},
			},
		},
		{
			Path: "type",
			Pattern: Pattern{
				Type:  PatternEquals,
				Value: "test",
			},
		},
	}

	data1, err := json.Marshal(r)
	assert.NilError(t, err)
	log.Println(string(data1))

	var u ConditionSet
	assert.NilError(t, json.Unmarshal(data1, &u))

	assert.Check(t, len(r) == len(u))

	data2, err := json.Marshal(u)
	assert.NilError(t, err)

	assert.Check(t, string(data1) == string(data2))
}

func TestConditions_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"production": {"equals": true}}`)

	var c ConditionSet
	err := json.Unmarshal(data, &c)
	assert.NilError(t, err)
	assert.Check(t, len(c) == 1)
	assert.Check(t, c[0].Path == "production")
	assert.Check(t, c[0].Pattern.Type == PatternEquals)
	assert.Check(t, c[0].Pattern.Value == "true")
}
