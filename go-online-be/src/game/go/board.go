package gogame

import "fmt"

const (
	SMALL_BOARD_SIZE  = 9
	MEDIUM_BOARD_SIZE = 13
	LARGE_BOARD_SIZE  = 19
)

type Square struct {
	stone uint8
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
		return 0, fmt.Errorf("invalid stone val %d", stone)
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

func (b *Board) isInBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Size && y < b.Size
}

func (b *Board) isEmpty(x, y int) bool {
	return b.Squares[x][y].stone == Empty
}
