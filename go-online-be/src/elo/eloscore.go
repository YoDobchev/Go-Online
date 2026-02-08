package elo

import (
	"fmt"
	"math"
)

const (
	FirstDanElo = 2100
	RankEloStep = 100
	DevFactor   = 20
)

func CalculateElo(playerA, playerB int, result float64) (int, int) {
	estimatedA := 1 / (1 + math.Pow(10, float64(playerA-playerB)/400))
	estimatedB := 1 - estimatedA

	newEloA := playerA + int(float64(DevFactor)*(result-estimatedA))
	newEloB := playerB + int(float64(DevFactor)*((1-result)-estimatedB))

	return newEloA, newEloB
}

func GetRank(elo int) (string, int) {
	rankIndex := elo - FirstDanElo

	if rankIndex < 0 {
		return fmt.Sprintf("%v kyu", -math.Floor(float64(rankIndex/RankEloStep))), rankIndex
	}
	return fmt.Sprintf("%v dan", math.Floor(float64(rankIndex/RankEloStep)+1)), rankIndex
}
