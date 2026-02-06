package gogame

import (
	"encoding/json"
	"fmt"
)

const (
	SMALL_BOARD_SIZE  = 9
	MEDIUM_BOARD_SIZE = 13
	LARGE_BOARD_SIZE  = 19
)

type Square struct {
	stone uint8
}

type Stone struct {
	x     int
	y     int
	color uint8
}

type Board struct {
	Squares [][]Square
	Size    int
}

func NewBoard(boardSize int) (*Board, error) {
	switch boardSize {
	case SMALL_BOARD_SIZE, MEDIUM_BOARD_SIZE, LARGE_BOARD_SIZE:
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedBoardSize, boardSize)
	}
	b := new(Board)
	b.Squares = make([][]Square, boardSize)
	b.Size = boardSize

	for i := range b.Squares {
		b.Squares[i] = make([]Square, boardSize)
	}
	return b, nil
}

func (b *Board) PrintBoard() {
	for i := range b.Squares {
		for j := range b.Squares[i] {
			r, err := stoneToEmoji(b.Squares[i][j].stone)
			if err != nil {
				r = '?'
			}
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
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
		return 0, fmt.Errorf("Invalid stone val %d", stone)
	}
}

func (b *Board) AddStone(x, y int, color uint8) error {
	if !b.isInBounds(x, y) {
		return fmt.Errorf("%w, x: %d, y: %d", ErrInvalidCoordinates, x, y)
	}
	if !b.isEmpty(x, y) {
		return fmt.Errorf("%w, x: %d, y: %d", ErrStoneAlreadyPlaced, x, y)
	}

	b.Squares[x][y].stone = color
	return nil
}

func (b *Board) RemoveStone(x, y int) {
	b.Squares[x][y].stone = Empty
}

func (b *Board) RemoveChain(chain []Stone) {
	for _, s := range chain {
		b.RemoveStone(s.x, s.y)
	}
}

func (b *Board) RemoveChains(chains [][]Stone) {
	for _, c := range chains {
		b.RemoveChain(c)
	}
}

func (b *Board) isInBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Size && y < b.Size
}

func (b *Board) isEmpty(x, y int) bool {
	return b.Squares[x][y].stone == Empty
}

func (b *Board) GetSquare(x, y int) (Square, error) {
	if !b.isInBounds(x, y) {
		return Square{}, fmt.Errorf("%w, x: %d, y: %d", ErrInvalidCoordinates, x, y)
	}

	return b.Squares[x][y], nil
}

func (b *Board) getNeighbors(x, y int) []Stone {
	neighbors := []Stone{}
	neighbors = b.appendStone(x-1, y, neighbors)
	neighbors = b.appendStone(x+1, y, neighbors)
	neighbors = b.appendStone(x, y-1, neighbors)
	neighbors = b.appendStone(x, y+1, neighbors)
	return neighbors
}

func (b *Board) appendStone(x, y int, stones []Stone) []Stone {
	square, err := b.GetSquare(x, y)
	if err == nil {
		stones = append(stones, Stone{x, y, square.stone})
	}
	return stones
}

func (b *Board) MarshalJSON() ([]byte, error) {
	grid := make([][]int, b.Size)

	for i := 0; i < b.Size; i++ {
		grid[i] = make([]int, b.Size)
		for j := 0; j < b.Size; j++ {
			grid[i][j] = int(b.Squares[i][j].stone)
		}
	}

	return json.Marshal(struct {
		Size    int     `json:"size"`
		Squares [][]int `json:"squares"`
	}{
		Size:    b.Size,
		Squares: grid,
	})
}

func (b *Board) UnmarshalJSON(data []byte) error {
	var aux struct {
		Size    int     `json:"size"`
		Squares [][]int `json:"squares"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	nb, err := NewBoard(aux.Size)
	if err != nil {
		return err
	}
	for x := 0; x < aux.Size; x++ {
		for y := 0; y < aux.Size; y++ {
			nb.Squares[x][y].stone = uint8(aux.Squares[x][y])
		}
	}
	*b = *nb
	return nil
}
