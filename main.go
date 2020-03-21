package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/buger/goterm"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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
			case !person.isAlive:
				color = goterm.YELLOW
			}
			_, err := goterm.Print(goterm.Color("*", color))
			if err != nil {
				log.Fatalf("error printing on screen: %v", err)
			}
		}
		goterm.Flush()
		currentDay++
		movePopulation(population, chanceOfTransmission, chanceOfDeath, maxX, maxY, currentDay, daysUntilDeath)
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
	data.AddColumn("Infected Count")
	for i, hi := range history {
		data.AddRow(float64(i), float64(hi.infectedCount))
	}
	goterm.Println(chart.Draw(data))
	goterm.Println("Total days lasted: ", len(history))
	var peopleInfected, peopleDead int
	for _, p := range population {
		if !p.isAlive {
			peopleDead++
		}
		if p.infectedAt != 0 {
			peopleInfected++
		}
	}
	goterm.Println("People infected: ", peopleInfected)
	goterm.Println("People died: ", peopleDead)
	goterm.Flush()
}

type historyItem struct {
	infectedCount  int
	healthyCount   int
	recoveredCount int
	deadCount      int
}

type position struct {
	x int
	y int
}

type person struct {
	position      position
	infectedAt    int
	isIsolated    bool
	isAlive       bool
	daysToRecover int
}

func (p person) isInfected(currentDay int) bool {
	return p.isAlive &&
		p.infectedAt != 0 &&
		currentDay-p.infectedAt < p.daysToRecover
}

func (p person) isImmune(currentDay int) bool {
	return p.isAlive && p.infectedAt != 0 && currentDay-p.infectedAt > p.daysToRecover
}

func (p person) daysInfected(currentDay int) int {
	return currentDay - p.infectedAt
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
					isAlive:       true,
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

func movePopulation(population []person, chanceOfTransmission, chanceOfDeath float64, maxX, maxY, currentDay, daysUntilDeath int) {
	for i := range population {
		if population[i].isIsolated || !population[i].isAlive {
			continue
		}
		if population[i].isInfected(currentDay) &&
			wouldOccurWithChance(chanceOfDeath) &&
			population[i].daysInfected(currentDay) == daysUntilDeath {
			population[i].isAlive = false
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
			if !p.isAlive {
				continue
			}
			if population[i].isInfected(currentDay) &&
				!p.isInfected(currentDay) &&
				!p.isImmune(currentDay) {
				if wouldOccurWithChance(chanceOfTransmission) {
					p.infectedAt = currentDay
				}
				continue
			}
			if p.isInfected(currentDay) && !population[i].isInfected(currentDay) &&
				!population[i].isImmune(currentDay) {
				if wouldOccurWithChance(chanceOfTransmission) {
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
		if !p.isAlive {
			historyItem.deadCount++
		}
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
