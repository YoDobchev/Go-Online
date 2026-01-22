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

	hash uint64
	seen map[uint64]struct{} // for superko

	chains *ChainsMap

	rules string
}

var (
	ErrUnsupportedBoardSize = errors.New("Unsupported board size")
	ErrInvalidCoordinates   = errors.New("Input coordinates are out of board bounds")
	ErrStoneAlreadyPlaced   = errors.New("Stone is already placed there")
	ErrIllegalKoMove        = errors.New("Illegal Ko move (Superko rule)")
)

func NewGame(boardSize int, players [2]string) (*Game, error) {
	g := new(Game)
	board, err := NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	g.Board = board
	g.chains = NewChainsMap(boardSize)

	g.Players = players
	g.CurrectTurn = Black

	zobristTable := NewZobristTable(boardSize)
	g.hash = zobristTable.HashBoard(g.Board.Squares, true)

	return g, nil
}

func switchTurn(g *Game) {
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

	err := g.Board.AddStone(x, y, g.CurrectTurn)
	if err != nil {
		return err
	}

	switchTurn(g)
	return nil
}

func (g *Game) Print() {
	g.Board.PrintBoard()
}
