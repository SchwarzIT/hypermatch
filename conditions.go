package hypermatch

import (
	"encoding/json"
	"strings"
)

type ConditionSet []Condition

func (c ConditionSet) MarshalJSON() ([]byte, error) {
	// CAUTION: this must not be a pointer-receiver!
	data := make(map[string]Pattern, len(c))
	for _, cc := range c {
		data[cc.Path] = cc.Pattern
	}

	return json.Marshal(data)
}

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

type Condition struct {
	Path    string  `json:"path"`
	Pattern Pattern `json:"pattern"`
}

func (c Condition) MarshalJSON() ([]byte, error) {
	// CAUTION: this must not be a pointer-receiver!
	return json.Marshal(map[string]Pattern{
		c.Path: c.Pattern,
	})
}

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

type Pattern struct {
	Type  PatternType `json:"type"`
	Value string      `json:"value,omitempty"`
	Sub   []Pattern   `json:"sub,omitempty"`
}

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
			d := string(v)
			if strings.HasPrefix(d, "\"") {
				d = d[1:]
			}
			if strings.HasSuffix(d, "\"") {
				d = d[:len(d)-1]
			}
			p.Value = d
		}
		break
	}
	return nil
}
