package hypermatch

import (
	"encoding/json"
)

// ConditionSet represents a rule and consists of one or more items of type Condition
type ConditionSet []Condition

// MarshalJSON marshals a ConditionSet into an easy readable JSON object
func (c ConditionSet) MarshalJSON() ([]byte, error) {
	// CAUTION: this must not be a pointer-receiver!
	data := make(map[string]Pattern, len(c))
	for _, cc := range c {
		data[cc.Path] = cc.Pattern
	}

	return json.Marshal(data)
}

// UnmarshalJSON unmarshal the JSON back to a ConditionSet
func (c *ConditionSet) UnmarshalJSON(data []byte) error {
	// CAUTION: this must be a pointer-receiver!
	var r map[string]Pattern
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	var cs ConditionSet
	for k, v := range r {
		cs = append(cs, Condition{
			Path:    k,
			Pattern: v,
		})
	}
	*c = cs
	return nil
}

// Condition represents a single condition inside a ConditionSet. It defines a Path (=reference to property in an event) and a Pattern to check against the value.
type Condition struct {
	Path    string  `json:"path"`
	Pattern Pattern `json:"pattern"`
}

// MarshalJSON marshals a Condition into an easy readable JSON object
func (c Condition) MarshalJSON() ([]byte, error) {
	// CAUTION: this must not be a pointer-receiver!
	return json.Marshal(map[string]Pattern{
		c.Path: c.Pattern,
	})
}

// UnmarshalJSON unmarshal the JSON back to a Condition
func (c *Condition) UnmarshalJSON(data []byte) error {
	// CAUTION: this must be a pointer-receiver!
	var r map[string]Pattern
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	for k, v := range r {
		c.Path = k
		c.Pattern = v
		break
	}
	return nil
}

// Pattern defines how a value should be compared. It consists of a Type and either a Value or Sub-patterns depending on the used Type.
type Pattern struct {
	Type  PatternType `json:"type"`
	Value string      `json:"value,omitempty"`
	Sub   []Pattern   `json:"sub,omitempty"`
}

// MarshalJSON marshals a Pattern into an easy readable JSON object
func (p Pattern) MarshalJSON() ([]byte, error) {
	// CAUTION: this must not be a pointer-receiver!

	if len(p.Sub) > 0 {
		return json.Marshal(map[string][]Pattern{
			p.Type.String(): p.Sub,
		})
	} else {
		return json.Marshal(map[string]string{
			p.Type.String(): p.Value,
		})
	}
}

// UnmarshalJSON unmarshal the JSON back to a Pattern
func (p *Pattern) UnmarshalJSON(data []byte) error {
	// CAUTION: this must be a pointer-receiver!

	var r map[string]json.RawMessage
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	for k, v := range r {
		p.Type = PatternTypeFromString(k)
		switch p.Type {
		case PatternAnythingBut, PatternAnyOf, PatternAllOf:
			var ps []Pattern
			if err := json.Unmarshal(v, &ps); err != nil {
				return err
			}
			p.Sub = ps
		default:
			var d string
			if err := json.Unmarshal(v, &d); err != nil {
				return err
			}
			p.Value = d
		}
		break
	}
	return nil
}
