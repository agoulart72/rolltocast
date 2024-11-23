package app

import (
	"math/rand"
)

type Strategies struct {
	MaxFirst           bool
	NoBacklashOnFail   bool
	RemoveCurrentLevel bool
}

type Stats struct {
	Level            int
	Strategies       Strategies
	NumberOfRuns     int
	SpellsPerRun     int
	NumberOfSuccess  []int
	NumberOfFailures []int
	MaxSuccess       []int
	MaxFailures      []int
	SuccessTimes     [][]int
	FailureTimes     [][]int
}

func RunSorcerer(level int, strategy Strategies, runs int, spellsPerRun int) Stats {

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

		totalCastingsOnRun := 0

		for mySorcerer.CurrentSpellRank > 0 && totalCastingsOnRun <= spellsPerRun {

			rank := rand.Intn(mySorcerer.CurrentSpellRank) + 1 // 0 -> max

			if strategy.MaxFirst {
				rank = mySorcerer.CurrentSpellRank
			}
			if strategy.RemoveCurrentLevel && !mySorcerer.SpellRanks[rank-1] {
				for i := range mySorcerer.SpellRanks {
					if mySorcerer.SpellRanks[i] {
						rank = i + 1
					}
				}
			}

			s, _ := mySorcerer.Cast(rank, strategy)
			if s {
				runSuccess[rank-1] += 1
			} else {
				runFailures[rank-1] += 1
			}

			totalCastingsOnRun++
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
		SpellsPerRun:     spellsPerRun,
		NumberOfSuccess:  numberOfSuccess,
		NumberOfFailures: numberOfFailures,
		MaxSuccess:       maxSuccess,
		MaxFailures:      maxFailures,
		SuccessTimes:     successTimes,
		FailureTimes:     failureTimes,
	}
}
