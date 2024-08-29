package main

import (
	"benchmark/candidates"
	"context"
	"log"
	"time"
)

const (
	numberOfRules      = 100000
	eventCheckDuration = 5 * time.Second
)

type Candidate interface {
	Name() string
	AddRule(number int, modulo int)
	Match(number int, modulo int) int
}

func main() {
	cs := []Candidate{candidates.NewHypermatch(), candidates.NewHypermatchJson(), candidates.NewQuamina()}

	for _, c := range cs {
		log.Printf("---Starting with %s\n", c.Name())
		beforeAddingRules := time.Now()

		for i := 0; i < numberOfRules; i++ {
			c.AddRule(i, numberOfRules/10)
		}
		log.Printf("adding %d rules took %.5fs\n", numberOfRules, time.Since(beforeAddingRules).Seconds())
		runEvents(c)
	}
}

func runEvents(c Candidate) {
	numberOfEvents := 0
	numberOfMatches := 0
	beforeCheckingEvents := time.Now()
	lastPrint := beforeCheckingEvents
	ctx, cancel := context.WithTimeout(context.Background(), eventCheckDuration)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			numberOfMatches += c.Match(numberOfMatches, numberOfRules/10)
			numberOfEvents += 1

			if time.Since(lastPrint).Seconds() >= 1 {
				log.Printf("processed %d events with %d matches in %.5fs -> %.5f evt/s\n", numberOfEvents, numberOfMatches, time.Since(beforeCheckingEvents).Seconds(), float64(numberOfEvents)/time.Since(beforeCheckingEvents).Seconds())
				lastPrint = time.Now()
			}
		}
	}
}
