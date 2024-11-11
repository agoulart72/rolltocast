package app

import "math/rand"

type Strategies struct {
	MaxFirst bool
}

type Stats struct {
	Level            int
	Strategies       Strategies
	NumberOfRuns     int
	NumberOfSuccess  []int
	NumberOfFailures []int
	MaxSuccess       []int
	MaxFailures      []int
	SuccessTimes     [][]int
	FailureTimes     [][]int
}

func RunSorcerer(level int, strategy Strategies, runs int) Stats {

	mySorcerer := &Sorcerer{
		Level: level,
	}

	numberOfSuccess := make([]int, 9)
	numberOfFailures := make([]int, 9)
	maxSuccess := make([]int, 9)
	maxFailures := make([]int, 9)
	successTimes := make([][]int, 9)
	failureTimes := make([][]int, 9)
	for i := 0; i < 9; i++ {
		successTimes[i] = make([]int, 10)
		failureTimes[i] = make([]int, 10)
	}

	for aRun := 0; aRun < runs; aRun++ {
		mySorcerer.LongRest()
		runSuccess := make([]int, 9)
		runFailures := make([]int, 9)

		for mySorcerer.CurrentSpellRank > 0 {

			rank := rand.Intn(mySorcerer.CurrentSpellRank) + 1 // 0 -> max

			if strategy.MaxFirst {
				rank = mySorcerer.CurrentSpellRank
			}

			s, _ := mySorcerer.Cast(rank)
			if s {
				runSuccess[rank-1] += 1
			} else {
				runFailures[rank-1] += 1
			}

		}

		for i := 0; i < 9; i++ {
			numberOfSuccess[i] += runSuccess[i]
			numberOfFailures[i] += runFailures[i]
			if runSuccess[i] > maxSuccess[i] {
				maxSuccess[i] = runSuccess[i]
			}
			if runFailures[i] > maxFailures[i] {
				maxFailures[i] = runFailures[i]
			}
			// Adjust size of success and failure arrays
			if runSuccess[i] >= cap(successTimes[i]) {
				newTimes := make([]int, runSuccess[i]+1)
				copy(newTimes, successTimes[i])
				successTimes[i] = newTimes
			}
			if runFailures[i] >= cap(failureTimes[i]) {
				newTimes := make([]int, runFailures[i]+1)
				copy(newTimes, failureTimes[i])
				failureTimes[i] = newTimes
			}
			successTimes[i][runSuccess[i]] += 1
			failureTimes[i][runFailures[i]] += 1
		}
	}

	return Stats{
		Level:            level,
		Strategies:       strategy,
		NumberOfRuns:     runs,
		NumberOfSuccess:  numberOfSuccess,
		NumberOfFailures: numberOfFailures,
		MaxSuccess:       maxSuccess,
		MaxFailures:      maxFailures,
		SuccessTimes:     successTimes,
		FailureTimes:     failureTimes,
	}
}
