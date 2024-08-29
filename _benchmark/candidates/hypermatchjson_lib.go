package candidates

import (
	"encoding/json"
	"fmt"
	"github.com/SchwarzIT/hypermatch"
)

type HypermatchJson struct {
	h *hypermatch.HyperMatch
}

func NewHypermatchJson() *HypermatchJson {
	return &HypermatchJson{h: hypermatch.NewHyperMatch()}
}

func (h *HypermatchJson) Name() string {
	return "hypermatch-json"
}

func (h *HypermatchJson) AddRule(number int, modulo int) {
	jsonStr := fmt.Sprintf(`
		{
			"name": {"wildcard": "*-myapp-*"},
			"env": {"equals": "prod"},
			"number": {"equals": "%d"},
			"tags": {"allOf": [{"equals": "tag1"}, {"equals": "tag2"}]},
			"region": {"anythingBut": [{"equals": "moon"}]},
			"type": {"anyOf": [{"equals": "app"}, {"equals": "database"}]}
		}
	`, number%modulo)
	var conditionSet hypermatch.ConditionSet
	if err := json.Unmarshal([]byte(jsonStr), &conditionSet); err != nil {
		panic(err)
	}
	if err := h.h.AddRule(hypermatch.RuleIdentifier(number), conditionSet); err != nil {
		panic(err)
	}
}

func (h *HypermatchJson) Match(number int, modulo int) int {
	eventStr := fmt.Sprintf(`
		[
			{"Path": "name", "Values": ["app-myapp-%d"]},
			{"Path": "env", "Values": ["prod"]},
			{"Path": "number", "Values": ["%d"]},
			{"Path": "tags", "Values": ["tag1", "tag2"]},
			{"Path": "region", "Values": ["earth"]},
			{"Path": "type", "Values": ["app"]}
		]
	`, number, number%modulo)

	var properties []hypermatch.Property
	if err := json.Unmarshal([]byte(eventStr), &properties); err != nil {
		panic(err)
	}

	matches := h.h.Match(properties)
	return len(matches)
}
