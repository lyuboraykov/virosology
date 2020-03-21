package main

import "time"

const (
	// populationSize is the size of the initially tested population
	populationSize = 400
	// initialInfectedCount is the count of people infected on day 1
	initialInfectedCount = 1
	// transmissionCoeff is the chance for a disease transmission to occur when two
	// people have contact. 0 to 1.
	transmissionCoeff = 0.8
	// isolationLevel is the portion of people who avoid contact with others. 0 to 1.
	isolationLevel = 0.1
	// daysToRecover is the time necessary for infected people to recover and become immune
	daysToRecover = 100
	// intervalBetweenFrames is the sleep time between renders
	intervalBetweenFrames = 300 * time.Millisecond
)
