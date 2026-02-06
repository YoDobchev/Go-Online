package gogame

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Move struct {
	Color uint8
	X, Y  int
}

const PASS = -1

const (
	Empty uint8 = iota
	White
	Black
)

const (
	GAME_WAITING_FOR_PLAYER = iota
	GAME_IN_PROGRESS
	GAME_ENDED
)

type Game struct {
	ID          string
	Players     [2]string
	CurrectTurn uint8
	passed      bool

	GameProgress uint8
	WhitePoints  int
	BlackPoints  int

	MoveNum int

	Board *Board

	hash    uint64
	seen    map[uint64]struct{} // for superko
	zobrist *Zobrist

	chains *ChainsMap

	rules string
}

var (
	ErrUnsupportedBoardSize = errors.New("Unsupported board size")
	ErrInvalidCoordinates   = errors.New("Input coordinates are out of board bounds")
	ErrStoneAlreadyPlaced   = errors.New("Stone is already placed there")
	ErrIllegalSelfCapture   = errors.New("Self capture is illegal if no opponent stones are captured")
	ErrIllegalKoMove        = errors.New("Illegal Ko move (Superko rule)")
)

func NewGame(boardSize int, creator string) (*Game, error) {
	g := new(Game)
	board, err := NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	g.Board = board
	g.chains = NewChainsMap(board)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.ID = fmt.Sprintf("%08d", rng.Intn(90000000)+10000000)
	g.GameProgress = GAME_WAITING_FOR_PLAYER
	GameInstances[g.ID] = g
	PlayerToGame[creator] = g

	g.Players[0] = creator
	g.CurrectTurn = Black

	g.zobrist = NewZobristTable(boardSize)
	g.hash = g.zobrist.HashBoard(g.Board.Squares, true)
	g.seen = make(map[uint64]struct{})

	saveGameToDB(g)
	saveSnapshotIfNeededToDB(g)

	return g, nil
}

func switchTurn(g *Game) {
	if g.CurrectTurn == Black {
		g.CurrectTurn = White
	} else {
		g.CurrectTurn = Black
	}
}

func (g *Game) Join(player string) error {
	if g.Players[1] != "" {
		return fmt.Errorf("game is full")
	}
	g.Players[1] = player
	PlayerToGame[player] = g
	saveGameToDB(g)
	return nil
}

func (g *Game) Leave(player string) error {
	if g.Players[0] != player && g.Players[1] != player {
		return nil
	}

	if g.Players[0] == player {
		g.Players[0] = ""
	}
	if g.Players[1] == player {
		g.Players[1] = ""
	}

	delete(PlayerToGame, player)

	if g.Players[0] == "" && g.Players[1] != "" {
		g.Players[0] = g.Players[1]
		g.Players[1] = ""
		PlayerToGame[g.Players[0]] = g
	}

	if g.Players[0] == "" && g.Players[1] == "" {
		delete(GameInstances, g.ID)
	}

	saveGameToDB(g)

	return nil
}

func (g *Game) PlayMove(player string, x, y int) error {
	if (g.CurrectTurn == Black && player != g.Players[0]) ||
		(g.CurrectTurn == White && player != g.Players[1]) {
		return fmt.Errorf("not your turn")
	}

	m := Move{
		Color: g.CurrectTurn,
		X:     x,
		Y:     y,
	}

	err := g.ApplyMove(m)
	if err != nil {
		return err
	}

	if g.GameProgress == GAME_WAITING_FOR_PLAYER {
		g.GameProgress = GAME_IN_PROGRESS
	}

	saveGameToDB(g)
	saveMoveToDB(g, m)
	saveSnapshotIfNeededToDB(g)

	g.MoveNum++

	return nil
}

func (g *Game) ApplyMove(m Move) error {
	if m.Color != g.CurrectTurn {
		return fmt.Errorf("wrong color for current turn")
	}

	if m.X == PASS {
		if g.passed {
			g.endGame()
			return nil
		}

		g.passed = true

		g.hash = g.zobrist.ToggleSide(g.hash)

		switchTurn(g)
		return nil
	}

	g.passed = false

	snapshot := g.chains.Snapshot()

	if err := g.Board.AddStone(m.X, m.Y, m.Color); err != nil {
		return err
	}

	sameColorCapturedChain, otherColorCapturedChains := g.chains.AddStone(m.X, m.Y, m.Color)

	if len(sameColorCapturedChain) > 0 && len(otherColorCapturedChains) == 0 {
		g.chains.Restore(snapshot)
		g.Board.RemoveStone(m.X, m.Y)
		return ErrIllegalSelfCapture
	}

	oldHash := g.hash
	g.hashMove(Stone{m.X, m.Y, m.Color}, otherColorCapturedChains)
	g.hash = g.zobrist.ToggleSide(g.hash)

	if _, found := g.seen[g.hash]; found {
		g.hash = oldHash
		g.chains.Restore(snapshot)
		g.Board.RemoveStone(m.X, m.Y)
		return ErrIllegalKoMove
	}

	g.Board.RemoveChains(otherColorCapturedChains)
	g.chains.applyCapture(sameColorCapturedChain, otherColorCapturedChains)

	switchTurn(g)
	g.seen[g.hash] = struct{}{}

	return nil
}

func (g *Game) endGame() {
	g.GameProgress = GAME_ENDED
	g.WhitePoints, g.BlackPoints = getPointsFromBoard(g.Board)
}

func (g *Game) hashMove(placedStone Stone, otherColorDeletedChains [][]Stone) {
	g.hash = g.zobrist.ToggleStone(g.hash, placedStone)
	for i := range otherColorDeletedChains {
		for _, stone := range otherColorDeletedChains[i] {
			g.hash = g.zobrist.ToggleStone(g.hash, stone)
		}
	}
}

func (g *Game) Print() {
	g.Board.PrintBoard()
}
