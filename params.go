package main

import "time"

const (
	// populationSize is the size of the initially tested population
	populationSize = 700
	// initialInfectedCount is the count of people infected on day 1
	initialInfectedCount = 10
	// chanceOfTransmission is the chance for a disease transmission to occur when two
	// people have contact. 0 to 1.
	chanceOfTransmission = 0.8
	// chanceOfDeath is the chance for an infected person to die. 0 to 1.
	chanceOfDeath = 0.05
	// daysUntilDeath is the number of days it would take to kill an infected person, if the person would die.
	daysUntilDeath = 20
	// isolationLevel is the portion of people who avoid contact with others. 0 to 1.
	isolationLevel = 0.1
	// daysToRecover is the time necessary for infected people to recover and become immune
	daysToRecover = 100
	// intervalBetweenFrames is the sleep time between renders
	intervalBetweenFrames = 100 * time.Millisecond
)
