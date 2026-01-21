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

type Square struct {
	stone uint8
}

const (
	BLACK_TURN = false
	WHITE_TURN = true
)

type Game struct {
	Players     [2]string
	CurrectTurn uint8

	Board [][]Square

	rules string

	groups map[int]*[]int
}

var (
	ErrUnsupportedBoardSize = errors.New("unsupported board size")
	ErrStoneAlreadyPlaced   = errors.New("Stone is already placed there")
)

func NewGame(boardSize int, players [2]string) (*Game, error) {
	switch boardSize {
	case 9, 13, 19:
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedBoardSize, boardSize)
	}

	g := new(Game)

	g.Board = make([][]Square, boardSize)
	for i := range g.Board {
		g.Board[i] = make([]Square, boardSize)
	}

	g.Players = players

	g.CurrectTurn = Black

	return g, nil
}

func switchTurn(g *Game) {
	if g.CurrectTurn == Black {
		g.CurrectTurn = White
	} else {
		g.CurrectTurn = Black
	}

}

func getNeighboursInfo(x, y int) {

}

func (g *Game) PlayMove(x, y int) error {
	if g.Board[x][y].stone != Empty {
		return fmt.Errorf("%w, x: %d, y; %d", ErrStoneAlreadyPlaced, x, y)
	}

	const PASS = -1
	if x == PASS {
		switchTurn(g)
		return nil
	}

	g.Board[x][y].stone = g.CurrectTurn

	switchTurn(g)
	return nil
}

func stoneToEmoji(stone uint8) (rune, error) {
	switch stone {
	case Empty:
		return '·', nil
	case White:
		return '○', nil
	case Black:
		return '●', nil
	default:
		return 0, fmt.Errorf("invalid stone val %d", stone)
	}
}

func (g *Game) PrintBoard() {
	for i := range g.Board {
		for j := range g.Board[i] {
			r, err := stoneToEmoji(g.Board[i][j].stone)
			if err != nil {
				r = '?'
			}
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
}
