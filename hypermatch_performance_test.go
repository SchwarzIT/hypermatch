package hypermatch

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestHyperMatchPerformance(t *testing.T) {
	h := NewHyperMatch()

	timer := time.Now()
	generateRules(h, 10000)
	log.Printf("generation took %dms", time.Since(timer)/time.Millisecond)

	data := []Property{
		{Path: "namespace", Values: []string{"monitoring-schwarz"}},
		{Path: "l1", Values: []string{"lidl"}},
		{Path: "name", Values: []string{"OS hdd"}},
		{Path: "priority", Values: []string{"P1"}},
		{Path: "instance", Values: []string{"de44"}},
		{Path: "hostgroups", Values: []string{"l", "ld", "de", "ss"}},
	}

	rounds := 1000

	count := 0
	index := 0

	timer = time.Now()

	go func() {
		for index < rounds-1 {
			<-time.After(1 * time.Second)
			log.Printf("%d alerts processed - %.1f alerts/s", index+1, float64(index)/float64(time.Since(timer)/time.Second))
		}
	}()

	for i := 0; i < rounds; i++ {
		count += len(h.Match(data))
		index = i
	}
	log.Printf("matching took %d ms, c: %d, c2: %d\n", time.Since(timer)/time.Millisecond, count, count/rounds)
}

func generateRules(h *HyperMatch, count int) {
	for i := 0; i < count; i++ {
		_ = h.AddRule(RuleIdentifier(i), []Condition{
			{Path: "l1", Pattern: Pattern{Type: PatternEquals, Value: "lidl"}},
			{Path: "name", Pattern: Pattern{Type: PatternPrefix, Value: "OS"}},
			{Path: "instance", Pattern: Pattern{Type: PatternEquals, Value: fmt.Sprintf("de%d", i%(count/100))}},
			{Path: "priority", Pattern: Pattern{Type: PatternAnyOf, Sub: []Pattern{{Type: PatternEquals, Value: "P1"}, {Type: PatternEquals, Value: "P2"}}}},
			{Path: "hostgroups", Pattern: Pattern{Type: PatternAllOf, Sub: []Pattern{{Type: PatternEquals, Value: "ld"}, {Type: PatternEquals, Value: "ss"}}}},
		})
	}

}
