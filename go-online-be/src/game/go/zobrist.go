package gogame

import "math/rand"

type Zobrist struct {
	Board      [][][]uint64 // [x][y][color]
	SideToMove uint64
}

func NewZobristTable(boardSize int) *Zobrist {
	var seed int64 = 0xC0FFEE
	r := rand.New(rand.NewSource(seed))

	board := make([][][]uint64, boardSize)
	for x := range boardSize {
		board[x] = make([][]uint64, boardSize)
		for y := range boardSize {
			board[x][y] = make([]uint64, 3)
			for c := White; c <= Black; c++ {
				board[x][y][c] = r.Uint64()
			}
		}
	}

	return &Zobrist{
		Board:      board,
		SideToMove: r.Uint64(),
	}
}

func (z *Zobrist) HashBoard(board [][]Square, blackToMove bool) uint64 {
	var h uint64
	for x := range board {
		for y := range board[x] {
			if board[x][y].stone != Empty {
				h ^= z.Board[x][y][board[x][y].stone]
			}
		}
	}
	if blackToMove {
		h ^= z.SideToMove
	}
	return h
}

func (z *Zobrist) ToggleStone(hash uint64, stone Stone) uint64 {
	return hash ^ z.Board[stone.x][stone.y][stone.color]
}

func (z *Zobrist) ToggleSide(hash uint64) uint64 {
	return hash ^ z.SideToMove
}
