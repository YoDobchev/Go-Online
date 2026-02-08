package gogame

import (
	"fmt"
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/YoDobchev/Go-Online/src/elo"
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

		var endedReason string
		if dbGame.GameEndedReason != nil {
			endedReason = *dbGame.GameEndedReason
		}

		g := &Game{
			ID:              dbGame.ID,
			Ranked:          dbGame.Ranked,
			Players:         [2]string{playerBlack, playerWhite},
			Board:           board,
			CurrectTurn:     dbGame.CurrentTurn,
			passed:          dbGame.Passed,
			WinnerIndex:     dbGame.WinnerIndex,
			GameEndedReason: endedReason,
			Komi:            dbGame.Komi,
			GameProgress:    dbGame.GameProgress,
			WhitePoints:     dbGame.WhitePoints,
			BlackPoints:     dbGame.BlackPoints,
			MoveNum:         dbGame.MoveNo,
			hash:            uint64(dbGame.CurrentHash),
			chains:          NewChainsMap(board),
			zobrist:         NewZobristTable(board.Size),
			seen:            make(map[uint64]struct{}),
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

	var endedReason *string
	if g.GameEndedReason != "" {
		endedReason = &g.GameEndedReason
	}

	dbGame := database.Game{
		ID:              g.ID,
		BoardSize:       g.Board.Size,
		Ranked:          g.Ranked,
		PlayerBlack:     playerBlack,
		PlayerWhite:     playerWhite,
		WinnerIndex:     g.WinnerIndex,
		GameEndedReason: endedReason,
		Komi:            g.Komi,
		CurrentTurn:     g.CurrectTurn,
		Passed:          g.passed,
		GameProgress:    g.GameProgress,
		WhitePoints:     g.WhitePoints,
		BlackPoints:     g.BlackPoints,
		MoveNo:          g.MoveNum,
		CurrentHash:     int64(g.hash),
		UpdatedAt:       time.Now(),
	}

	return database.DB.Save(&dbGame).Error
}

func updateElo(g *Game) error {
	if database.DB == nil {
		return nil
	}

	winnerName := g.Players[g.WinnerIndex]
	loserName := g.Players[1-g.WinnerIndex]

	winner, err := loadUserFromDB(winnerName)
	if err != nil {
		return err
	}

	loser, err := loadUserFromDB(loserName)
	if err != nil {
		return err
	}

	winnerElo, loserElo := elo.CalculateElo(winner.Elo, loser.Elo, 1)
	winner.Elo, loser.Elo = winnerElo, loserElo

	if err := database.DB.Save(winner).Error; err != nil {
		return err
	}
	if err := database.DB.Save(loser).Error; err != nil {
		return err
	}

	return nil
}

func loadUserFromDB(username string) (*database.User, error) {
	var user database.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
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

func LoadBlogFromDB(blogID int) (map[string]any, error) {
	var blog map[string]any
	err := database.DB.
		Model(&database.Blog{}).
		Where("id = ?", blogID).
		First(&blog).Error
	return blog, err
}

func CreateBlog(id int, authorID int, title string, content string) error {
	var author database.User
	if err := database.DB.First(&author, authorID).Error; err != nil {
		return err
	}

	blog := database.Blog{
		ID:          id,
		AuthorID:    authorID,
		AuthorName:  author.Username,
		Title:       title,
		BlogContent: content,
		PublishedAt: time.Now(),
	}

	return database.DB.Create(&blog).Error
}

func GetLeaderBoard() ([]map[string]any, error) {
	var leaderboard []map[string]any
	err := database.DB.
		Model(&database.User{}).
		Order("elo DESC").
		Limit(20).
		Find(&leaderboard).Error
	for _, user := range leaderboard {
		user["rank"], _ = elo.GetRank(user["elo"].(int))
	}
	return leaderboard, err
}

func GetBlogReplies(blogID int) ([]map[string]any, error) {
	var replies []map[string]any
	err := database.DB.
		Model(&database.BlogReply{}).Where("blog_id = ?", blogID).Find(&replies).Error
	return replies, err
}

func CreateBlogReply(blogID int, authorID int, content string) error {
	var author database.User
	if err := database.DB.First(&author, authorID).Error; err != nil {
		return err
	}

	var blog database.Blog
	if err := database.DB.First(&blog, blogID).Error; err != nil {
		return err
	}

	reply := database.BlogReply{
		BlogID:       blogID,
		AuthorID:     authorID,
		AuthorName:   author.Username,
		ReplyContent: content,
		CreatedAt:    time.Now(),
	}

	return database.DB.Create(&reply).Error
}

func DeleteBlogReply(authorID int, replyID int) error {
	var author database.User
	if err := database.DB.First(&author, authorID).Error; err != nil {
		return err
	}

	var reply database.BlogReply
	if err := database.DB.First(&reply, replyID).Error; err != nil {
		return err
	}

	if author.Role != "admin" && author.ID != reply.AuthorID {
		return fmt.Errorf("user is not authorized to delete this blog reply")
	}

	return database.DB.Delete(&reply).Error
}

func SaveReportToDB(gameID string, username string) error {
	if database.DB == nil {
		return nil
	}

	report := database.Report{
		GameID:   gameID,
		Username: username,
	}
	return database.DB.Create(&report).Error
}

func LoadAllReportsFromDB() ([]map[string]any, error) {
	var reports []map[string]any

	err := database.DB.
		Model(&database.Report{}).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, err
}

func DeleteReportFromDB(reportID string) error {
	if database.DB == nil {
		return nil
	}

	return database.DB.Delete(&database.Report{}, reportID).Error
}
