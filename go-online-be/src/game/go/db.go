package gogame

import "github.com/YoDobchev/Go-Online/src/database"

func saveMoveToDB(g *Game, m Move) error {
	move := database.GameMove{
		GameID: g.ID,
		MoveNo: g.MoveNum,
		Color:  m.Color,
		X:      m.X,
		Y:      m.Y,
	}

	return database.DB.Create(&move).Error
}

func saveGameToDB(g *Game) error {
	var playerBlack, playerWhite *string
	if g.Players[0] != "" {
		playerBlack = &g.Players[0]
	}
	if g.Players[1] != "" {
		playerWhite = &g.Players[1]
	}

	dbGame := database.Game{
		ID:           g.ID,
		BoardSize:    g.Board.Size,
		PlayerBlack:  playerBlack,
		PlayerWhite:  playerWhite,
		CurrentTurn:  g.CurrectTurn,
		Passed:       g.passed,
		GameProgress: g.GameProgress,
		WhitePoints:  g.WhitePoints,
		BlackPoints:  g.BlackPoints,
	}

	return database.DB.Save(&dbGame).Error
}
