package gogame

import (
	"errors"
	"fmt"
)

const (
	Empty uint8 = iota
	White
	Black
)

type Groups struct {
	Map        [][]uint64 // [x][y][group]
	GroupCount uint64
	Liberties  map[uint64]uint64 // [group] -> num of liberties
}

type Game struct {
	Players     [2]string
	CurrectTurn uint8

	Board *Board

	hash uint64
	seen map[uint64]struct{} // for superko

	groups Groups

	rules string
}

var (
	ErrUnsupportedBoardSize = errors.New("unsupported board size")
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
	if !g.Board.isEmpty(x, y) {
		return fmt.Errorf("%w, x: %d, y; %d", ErrStoneAlreadyPlaced, x, y)
	}

	const PASS = -1
	if x == PASS {
		switchTurn(g)
		return nil
	}

	g.Board.Squares[x][y].stone = g.CurrectTurn

	switchTurn(g)
	return nil
}

func (g *Game) Print() {
	g.Board.PrintBoard()
}
