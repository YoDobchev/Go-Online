package gogame

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const PASS = -1

const (
	Empty uint8 = iota
	White
	Black
)

type Game struct {
	ID          string
	Players     [2]string
	CurrectTurn uint8

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
	GameInstances[g.ID] = g
	PlayerToGame[creator] = g

	g.Players[0] = creator
	g.CurrectTurn = Black

	g.zobrist = NewZobristTable(boardSize)
	g.hash = g.zobrist.HashBoard(g.Board.Squares, true)
	g.seen = make(map[uint64]struct{})

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

	return nil
}

func (g *Game) PlayMove(player string, x, y int) error {
	if (g.CurrectTurn == Black && player != g.Players[0]) ||
		(g.CurrectTurn == White && player != g.Players[1]) {
		return fmt.Errorf("not your turn")
	}

	if x == PASS {
		switchTurn(g)
		return nil
	}

	snapshot := g.chains.Snapshot()
	err := g.Board.AddStone(x, y, g.CurrectTurn)
	if err != nil {
		return err
	}
	sameColorCapturedChain, otherColorCapturedChains := g.chains.AddStone(x, y, g.CurrectTurn)

	if len(sameColorCapturedChain) > 0 && len(otherColorCapturedChains) == 0 {
		//Self capture
		g.chains.Restore(snapshot)
		g.Board.RemoveStone(x, y)
		return ErrIllegalSelfCapture
	}
	oldHash := g.hash
	g.hashMove(Stone{x, y, g.CurrectTurn}, otherColorCapturedChains)
	g.hash = g.zobrist.ToggleSide(g.hash)
	_, found := g.seen[g.hash]
	if found {
		//Ko move
		g.hash = oldHash
		g.chains.Restore(snapshot)
		g.Board.RemoveStone(x, y)
		return ErrIllegalKoMove
	} else {
		//Valid move
		g.Board.RemoveChains(otherColorCapturedChains)
		g.chains.applyCapture(sameColorCapturedChain, otherColorCapturedChains)
	}
	switchTurn(g)
	g.seen[g.hash] = struct{}{}
	fmt.Println(len(g.seen))
	return nil
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
