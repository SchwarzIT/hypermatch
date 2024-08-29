package candidates

import (
	"fmt"
	"log"
	"quamina.net/go/quamina"
)

type Quamina struct {
	q *quamina.Quamina
}

func NewQuamina() *Quamina {
	q, err := quamina.New(quamina.WithMediaType("application/json"))
	if err != nil {
		panic(err)
	}
	return &Quamina{q: q}
}

func (q *Quamina) Name() string {
	return "quamina"
}

func (q *Quamina) AddRule(number int, modulo int) {
	str := fmt.Sprintf(`
		{
			"name": [{"shellstyle": "*-myapp-*"}],
			"env": ["prod"],
			"nunmber": ["%d"],
			"tags": ["tag1", "tag2"],
			"region": [{"anything-but": ["moon"]}],
			"type": ["app", "database"]
		}
	`, number%modulo)
	err := q.q.AddPattern(number, str)
	if err != nil {
		log.Panicln(err)
	}
}

func (q *Quamina) Match(number int, modulo int) int {
	event := fmt.Sprintf(`
		{
			"name": "app-myapp-%d",
			"env": "prod",
			"number": "%d",
			"tags": ["tag1", "tag2"],
			"region": "earth",
			"type": "app"
		}
	`, number, number%modulo)
	r, _ := q.q.MatchesForEvent([]byte(event))
	return len(r)
}
