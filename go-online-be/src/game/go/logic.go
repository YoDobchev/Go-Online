package gogame

import (
	"errors"
)

const PASS = -1

const (
	Empty uint8 = iota
	White
	Black
)

type Game struct {
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

func NewGame(boardSize int, players [2]string) (*Game, error) {
	g := new(Game)
	board, err := NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	g.Board = board
	g.chains = NewChainsMap(board)

	g.Players = players
	g.CurrectTurn = Black

	g.zobrist = NewZobristTable(boardSize)
	g.hash = g.zobrist.HashBoard(g.Board.Squares, true)
	g.seen = make(map[uint64]struct{})

	return g, nil
}

func switchTurn(g *Game) {
	g.hash = g.zobrist.ToggleSide(g.hash)
	if g.CurrectTurn == Black {
		g.CurrectTurn = White
	} else {
		g.CurrectTurn = Black
	}
}

func (g *Game) PlayMove(x, y int) error {
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
	g.hashMove(Stone{x, y, g.CurrectTurn}, sameColorCapturedChain, otherColorCapturedChains)
	_, found := g.seen[g.hash]
	if found {
		//Ko move
		g.hash = oldHash
		g.chains.Restore(snapshot)
		g.Board.RemoveStone(x, y)
		return ErrIllegalKoMove
	} else {
		//Valid move
		g.Board.RemoveChain(sameColorCapturedChain)
		g.Board.RemoveChains(otherColorCapturedChains)
		g.chains.applyCapture(sameColorCapturedChain, otherColorCapturedChains)
	}

	switchTurn(g)
	g.seen[g.hash] = struct{}{}
	return nil
}

func (g *Game) hashMove(placedStone Stone, sameColorDeletedChain []Stone, otherColorDeletedChains [][]Stone) {
	g.hash = g.zobrist.ToggleStone(g.hash, placedStone)
	for _, stone := range sameColorDeletedChain {
		g.hash = g.zobrist.ToggleStone(g.hash, stone)
	}
	for i := range otherColorDeletedChains {
		for _, stone := range otherColorDeletedChains[i] {
			g.hash = g.zobrist.ToggleStone(g.hash, stone)
		}
	}
}

func (g *Game) Print() {
	g.Board.PrintBoard()
	//g.chains.PrintGroups()
}
