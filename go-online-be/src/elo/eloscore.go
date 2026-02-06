package elo

import (
	"fmt"
	"math"
)

const (
	FirstDanElo = 2100
	RankEloStep = 100
)

type Profile struct {
	UserID    string
	Username  string
	Elo       int
	DevFactor int
}

func CalculateElo(playerA, playerB Profile, result float64) (int, int) {
	estimatedA := 1 / (1 + math.Pow(10, float64(playerA.Elo-playerB.Elo)/400))
	estimatedB := 1 - estimatedA

	newEloA := playerA.Elo + int(float64(playerA.DevFactor)*(result-estimatedA))
	newEloB := playerB.Elo + int(float64(playerB.DevFactor)*((1-result)-estimatedB))

	return newEloA, newEloB
}

func GetRank(elo int) (string, int) {
	rankIndex := elo - FirstDanElo

	if rankIndex < 0 {
		return fmt.Sprintf("%v kyu", -math.Floor(float64(rankIndex/RankEloStep))), rankIndex
	}
	return fmt.Sprintf("%v dan", math.Floor(float64(rankIndex/RankEloStep)+1)), rankIndex
}
