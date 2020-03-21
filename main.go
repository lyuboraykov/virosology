package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/buger/goterm"
)

func main() {
	currentDay := 1
	maxX := goterm.Width()
	maxY := goterm.Height()
	population := newPopulation(populationSize, initialInfectedCount, daysToRecover,
		isolationLevel, maxX, maxY)

	for {
		goterm.Clear()
		for _, person := range population {
			goterm.MoveCursor(person.position.x, person.position.y)
			color := goterm.WHITE
			switch {
			case person.isInfected(currentDay):
				color = goterm.RED
			case person.isImmune(currentDay):
				color = goterm.GREEN
			}
			_, err := goterm.Print(goterm.Color("*", color))
			if err != nil {
				log.Fatalf("error printing on screen: %v", err)
			}
		}
		goterm.Flush()
		currentDay++
		movePopulation(population, maxX, maxY, currentDay)
		time.Sleep(intervalBetweenFrames)
	}
}

type position struct {
	x int
	y int
}

type person struct {
	position      position
	infectedAt    int
	isIsolated    bool
	daysToRecover int
}

func (p person) isInfected(currentDay int) bool {
	return p.infectedAt != 0 && currentDay-p.infectedAt < p.daysToRecover
}

func (p person) isImmune(currentDay int) bool {
	return p.infectedAt != 0 && currentDay-p.infectedAt > p.daysToRecover
}

func newPopulation(
	populationSize,
	initialInfectedCount,
	daysToRecover int,
	isolationLevel float32,
	maxX,
	maxY int,
) []person {
	population := make([]person, populationSize)
	for i := 0; i < populationSize; i++ {
		for {
			x := rand.Intn(maxX) + 1
			y := rand.Intn(maxY) + 1
			candidatePosition := position{x, y}
			if _, taken := positionTaken(candidatePosition, population); !taken {
				var infectedAt int
				if i < initialInfectedCount {
					infectedAt = 1
				}
				population[i] = person{
					position:      candidatePosition,
					infectedAt:    infectedAt,
					isIsolated:    rand.Intn(100) < int(isolationLevel*100),
					daysToRecover: daysToRecover,
				}
				break
			}
		}
	}

	return population
}

func positionTaken(pos position, pop []person) (*person, bool) {
	for _, p := range pop {
		if p.position == pos {
			return &p, true
		}
	}
	return nil, false
}

func movePopulation(population []person, maxX, maxY, currentDay int) {
	for i := range population {
		if population[i].isIsolated {
			continue
		}
		xOrY := rand.Intn(100) < 50
		minusOrPlus := rand.Intn(100) < 50
		direction := 1
		if minusOrPlus {
			direction = -1
		}
		candidatePosition := population[i].position
		if xOrY {
			candidatePosition.x += direction
		} else {
			candidatePosition.y += direction
		}

		if candidatePosition.x < 1 || candidatePosition.y < 1 || candidatePosition.x > maxX || candidatePosition.y > maxY {
			continue
		}

		if p, taken := positionTaken(candidatePosition, population); taken {
			if population[i].isInfected(currentDay) &&
				!p.isInfected(currentDay) &&
				!p.isImmune(currentDay) {
				p.infectedAt = currentDay
				continue
			}
			if p.isInfected(currentDay) && !population[i].isInfected(currentDay) &&
				!population[i].isImmune(currentDay) {
				population[i].infectedAt = currentDay
				continue
			}
		}

		population[i].position = candidatePosition
	}
}
