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

	var history []historyItem

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
		movePopulation(population, transmissionCoeff, maxX, maxY, currentDay)
		hi := getHistoryItem(population, currentDay)
		history = append(history, hi)
		if hi.infectedCount == 0 {
			break
		}
		time.Sleep(intervalBetweenFrames)
	}
	goterm.Clear()
	goterm.MoveCursor(1, 1)
	goterm.Flush()
	chart := goterm.NewLineChart(maxX-10, maxY-10)
	data := new(goterm.DataTable)
	data.AddColumn("Time")
	data.AddColumn(goterm.Color("Infected Count", goterm.RED))
	data.AddColumn(goterm.Color("Recovered Count", goterm.GREEN))
	data.AddColumn(goterm.Color("Healthy Count", goterm.WHITE))
	for i, hi := range history {
		data.AddRow(float64(i), float64(hi.infectedCount),
			float64(hi.recoveredCount), float64(hi.healthyCount))
	}
	goterm.Println(chart.Draw(data))
	goterm.Println("Total days lasted: ", len(history))
	var peopleInfected int
	for _, p := range population {
		if p.infectedAt != 0 {
			peopleInfected++
		}
	}
	goterm.Println("People infected: ", peopleInfected)
	goterm.Flush()
}

type historyItem struct {
	infectedCount  int
	healthyCount   int
	recoveredCount int
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
	isolationLevel float64,
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
					isIsolated:    wouldOccurWithChance(isolationLevel),
					daysToRecover: daysToRecover,
				}
				break
			}
		}
	}

	return population
}

func positionTaken(pos position, population []person) (*person, bool) {
	for i := range population {
		if population[i].position == pos {
			return &population[i], true
		}
	}
	return nil, false
}

func movePopulation(population []person, transmissionCoeff float64, maxX, maxY, currentDay int) {
	for i := range population {
		if population[i].isIsolated {
			continue
		}
		xOrY := wouldOccurWithChance(0.5)
		minusOrPlus := wouldOccurWithChance(0.5)
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
				if wouldOccurWithChance(transmissionCoeff) {
					p.infectedAt = currentDay
				}
				continue
			}
			if p.isInfected(currentDay) && !population[i].isInfected(currentDay) &&
				!population[i].isImmune(currentDay) {
				if wouldOccurWithChance(transmissionCoeff) {
					population[i].infectedAt = currentDay
				}
				continue
			}
		}

		population[i].position = candidatePosition
	}
}

func getHistoryItem(population []person, currentDay int) historyItem {
	historyItem := historyItem{}
	for _, p := range population {
		if p.isInfected(currentDay) {
			historyItem.infectedCount++
		} else if p.isImmune(currentDay) {
			historyItem.recoveredCount++
		} else {
			historyItem.healthyCount++
		}
	}
	return historyItem
}

func wouldOccurWithChance(chance float64) bool {
	return rand.Intn(100) < int(chance*100)
}
