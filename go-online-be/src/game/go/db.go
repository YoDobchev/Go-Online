package gogame

import (
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
)

const (
	SnapshotInterval = 3
)

func LoadGamesFromDB() {
	var dbGames []database.Game
	if err := database.DB.Find(&dbGames).Error; err != nil {
		panic(err)
	}

	for _, dbGame := range dbGames {
		var playerBlack, playerWhite string

		if dbGame.PlayerBlack != nil {
			playerBlack = *dbGame.PlayerBlack
		}
		if dbGame.PlayerWhite != nil {
			playerWhite = *dbGame.PlayerWhite
		}
		board, err := GetBoardStateOnMoveNoFromDB(dbGame.ID, dbGame.MoveNo)
		if err != nil {
			panic(err)
		}
		g := &Game{
			ID:           dbGame.ID,
			Ranked:       dbGame.Ranked,
			Players:      [2]string{playerBlack, playerWhite},
			Board:        board,
			CurrectTurn:  dbGame.CurrentTurn,
			passed:       dbGame.Passed,
			WinnerIndex:  dbGame.WinnerIndex,
			Komi:         dbGame.Komi,
			GameProgress: dbGame.GameProgress,
			WhitePoints:  dbGame.WhitePoints,
			BlackPoints:  dbGame.BlackPoints,
			MoveNum:      dbGame.MoveNo + 1,
			hash:         uint64(dbGame.CurrentHash),
			chains:       NewChainsMap(board),
			zobrist:      NewZobristTable(board.Size),
			seen:         make(map[uint64]struct{}),
		}

		GameInstances[g.ID] = g
		if g.GameProgress != GAME_ENDED {
			if g.Players[0] != "" {
				PlayerToGame[g.Players[0]] = g
			}
			if g.Players[1] != "" {
				PlayerToGame[g.Players[1]] = g
			}
		}
	}
}

func constructGameFromSnapshot(snapshot database.GameSnapshot) (*Game, error) {
	board := &Board{}
	if err := board.UnmarshalJSON(snapshot.BoardJSON); err != nil {
		return nil, err
	}

	g := &Game{
		ID:          snapshot.GameID,
		Board:       board,
		CurrectTurn: snapshot.CurrentTurn,
		passed:      snapshot.Passed,
		hash:        uint64(snapshot.Hash),
		chains:      NewChainsMap(board),
		zobrist:     NewZobristTable(board.Size),
		seen:        make(map[uint64]struct{}),
	}

	return g, nil
}

func GetBoardStateOnMoveNoFromDB(gameID string, moveNo int) (*Board, error) {
	var snap database.GameSnapshot

	err := database.DB.
		Where("game_id = ? AND move_no <= ?", gameID, moveNo).
		Order("move_no DESC").
		First(&snap).Error
	if err != nil {
		return nil, err
	}

	g, err := constructGameFromSnapshot(snap)
	if err != nil {
		return nil, err
	}

	var moves []database.GameMove
	err = database.DB.
		Where("game_id = ? AND move_no > ? AND move_no <= ?", gameID, snap.MoveNo, moveNo).
		Order("move_no ASC").
		Find(&moves).Error
	if err != nil {
		return nil, err
	}

	for _, mv := range moves {
		if err := g.ApplyMove(Move{
			Color: mv.Color,
			X:     mv.X,
			Y:     mv.Y,
		}); err != nil {
			return nil, err
		}
	}

	return g.Board, nil
}

func saveMoveToDB(g *Game, m Move) error {
	if database.DB == nil {
		return nil
	}

	move := database.GameMove{
		GameID:        g.ID,
		MoveNo:        g.MoveNum,
		Color:         m.Color,
		X:             m.X,
		Y:             m.Y,
		ResultingHash: int64(g.hash),
	}

	return database.DB.Create(&move).Error
}

func saveGameToDB(g *Game) error {
	if database.DB == nil {
		return nil
	}

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
		Ranked:       g.Ranked,
		PlayerBlack:  playerBlack,
		PlayerWhite:  playerWhite,
		WinnerIndex:  g.WinnerIndex,
		Komi:         g.Komi,
		CurrentTurn:  g.CurrectTurn,
		Passed:       g.passed,
		GameProgress: g.GameProgress,
		WhitePoints:  g.WhitePoints,
		BlackPoints:  g.BlackPoints,
		MoveNo:       g.MoveNum,
		CurrentHash:  int64(g.hash),
		UpdatedAt:    time.Now(),
	}

	return database.DB.Save(&dbGame).Error
}

func saveSnapshotIfNeededToDB(g *Game) error {
	if database.DB == nil {
		return nil
	}

	if g.MoveNum%SnapshotInterval != 0 && g.GameProgress != GAME_ENDED {
		return nil
	}

	boardJSON, err := g.Board.MarshalJSON()
	if err != nil {
		return err
	}

	dbSnapshot := database.GameSnapshot{
		GameID: g.ID,
		MoveNo: g.MoveNum,

		BoardJSON:   boardJSON,
		CurrentTurn: g.CurrectTurn,
		Passed:      g.passed,
		Hash:        int64(g.hash),

		CreatedAt: time.Now(),
	}

	return database.DB.Save(&dbSnapshot).Error
}

func LoadAllBlogsFromDB() ([]map[string]any, error) {
	var blogs []map[string]any

	err := database.DB.
		Model(&database.Blog{}).
		Order("published_at DESC").
		Omit("blog_content").
		Find(&blogs).Error

	return blogs, err
}
