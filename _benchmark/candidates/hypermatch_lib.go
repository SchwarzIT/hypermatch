package candidates

import (
	"fmt"
	"github.com/SchwarzIT/hypermatch"
	"log"
)

type Hypermatch struct {
	h *hypermatch.HyperMatch
}

func NewHypermatch() *Hypermatch {
	return &Hypermatch{h: hypermatch.NewHyperMatch()}
}

func (h *Hypermatch) Name() string {
	return "hypermatch"
}

func (h *Hypermatch) AddRule(number int, modulo int) {
	err := h.h.AddRule(hypermatch.RuleIdentifier(number), hypermatch.ConditionSet{
		{Path: "name", Pattern: hypermatch.Pattern{Type: hypermatch.PatternWildcard, Value: "*-myapp-*"}},
		{Path: "env", Pattern: hypermatch.Pattern{Type: hypermatch.PatternEquals, Value: "prod"}},
		{Path: "number", Pattern: hypermatch.Pattern{Type: hypermatch.PatternEquals, Value: fmt.Sprintf("%d", number%modulo)}},
		{Path: "tags", Pattern: hypermatch.Pattern{
			Type: hypermatch.PatternAllOf, Sub: []hypermatch.Pattern{
				{Type: hypermatch.PatternEquals, Value: "tag1"},
				{Type: hypermatch.PatternEquals, Value: "tag2"},
			},
		}},
		{Path: "region", Pattern: hypermatch.Pattern{
			Type: hypermatch.PatternAnythingBut, Sub: []hypermatch.Pattern{
				{Type: hypermatch.PatternEquals, Value: "moon"},
			},
		}},
		{Path: "type", Pattern: hypermatch.Pattern{
			Type: hypermatch.PatternAnyOf, Sub: []hypermatch.Pattern{
				{Type: hypermatch.PatternEquals, Value: "app"},
				{Type: hypermatch.PatternEquals, Value: "database"},
			},
		}},
	})
	if err != nil {
		log.Panicln(err)
	}
}

func (h *Hypermatch) Match(number int, modulo int) int {
	event := []hypermatch.Property{
		{
			Path:   "name",
			Values: []string{fmt.Sprintf("app-myapp-%d", number)},
		},
		{
			Path:   "env",
			Values: []string{"prod"},
		},
		{
			Path:   "number",
			Values: []string{fmt.Sprintf("%d", number%modulo)},
		},
		{
			Path:   "tags",
			Values: []string{"tag1", "tag2"},
		},
		{
			Path:   "region",
			Values: []string{"earth"},
		},
		{
			Path:   "type",
			Values: []string{"app"},
		},
	}
	matches := h.h.Match(event)
	return len(matches)
}
