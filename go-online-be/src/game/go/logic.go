package gogame

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	clock "github.com/YoDobchev/Go-Online/src/game/clock"
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

type NewGameSettings struct {
	PlayAs    int
	BoardSize int
	Ranked    bool
	// Time      string `json:"time"`
}

type Game struct {
	ID          string
	Ranked      bool
	Players     [2]string
	WinnerIndex uint8
	CurrectTurn uint8
	passed      bool

	GameProgress    uint8
	WhitePoints     int
	BlackPoints     int
	GameEndedReason string

	done     chan struct{}
	doneOnce sync.Once

	Komi float32

	MoveNum int

	Board *Board

	hash    uint64
	seen    map[uint64]struct{} // for superko
	zobrist *Zobrist

	chains *ChainsMap

	rules string

	Clock      *clock.Clock
	Events     chan map[string]any
	MovePlayed chan struct{}
	mu         sync.Mutex
}

var (
	ErrUnsupportedBoardSize = errors.New("Unsupported board size")
	ErrInvalidCoordinates   = errors.New("Input coordinates are out of board bounds")
	ErrStoneAlreadyPlaced   = errors.New("Stone is already placed there")
	ErrIllegalSelfCapture   = errors.New("Self capture is illegal if no opponent stones are captured")
	ErrIllegalKoMove        = errors.New("Illegal Ko move (Superko rule)")
)

func NewGame(creator string, settings NewGameSettings) (*Game, error) {
	g := new(Game)
	board, err := NewBoard(settings.BoardSize)
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

	if settings.PlayAs == 2 {
		settings.PlayAs = rand.Intn(2)
	}

	g.Players[settings.PlayAs] = creator
	g.CurrectTurn = Black
	g.Ranked = settings.Ranked

	g.zobrist = NewZobristTable(settings.BoardSize)
	g.hash = g.zobrist.HashBoard(g.Board.Squares, true)
	g.seen = make(map[uint64]struct{})

	g.Events = make(chan map[string]any, 64)
	g.MovePlayed = make(chan struct{}, 16)
	g.done = make(chan struct{})

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
	if g.Players[0] == player || g.Players[1] == player {
		return nil
	}

	if g.GameProgress == GAME_ENDED {
		return fmt.Errorf("game has ended")
	}

	if g.Players[0] != "" && g.Players[1] != "" {
		return fmt.Errorf("game is full")
	}

	newPlayerIndex := 0
	if g.Players[1] == "" {
		newPlayerIndex = 1
	}

	g.Players[newPlayerIndex] = player

	if g.Players[0] != "" && g.Players[1] != "" {
		if g.GameProgress == GAME_WAITING_FOR_PLAYER {
			g.GameProgress = GAME_IN_PROGRESS
			g.startClock()
		}
	}

	PlayerToGame[player] = g
	saveGameToDB(g)
	return nil
}

func (g *Game) startClock() {
	go func() {
		timeFormat := clock.TimeFormat{
			MainTime:       300 * time.Second,
			ByoYomi:        0,
			ByoYomiPeriods: 0,
		}

		g.Clock = clock.NewClock(timeFormat)
		g.Clock.Start(Black)
		turn := Black

		tick := time.NewTicker(100 * time.Millisecond)
		update := time.NewTicker(1 * time.Second)
		defer tick.Stop()
		defer update.Stop()
		g.Events <- g.Clock.GetClockUpdate(Black)
		g.Events <- g.Clock.GetClockUpdate(White)
		for {
			select {
			case <-update.C:
				g.Events <- g.Clock.GetClockUpdate(turn)
			case <-tick.C:
				g.mu.Lock()
				turn = g.CurrectTurn
				ended := g.GameProgress == GAME_ENDED
				g.mu.Unlock()

				if ended {
					return
				}

				if g.Clock.OutOfTime(turn) {
					g.mu.Lock()
					loser := turn
					g.GameEndedReason = "timeout"
					g.GameProgress = GAME_ENDED

					if loser == Black {
						g.WinnerIndex = 1 // white wins
					} else {
						g.WinnerIndex = 0 // black wins
					}
					g.mu.Unlock()

					// g.Events <- map[string]any{
					// 	"type": "game_ended",
					// 	"data": map[string]any{
					// 		"white_points": g.WhitePoints,
					// 		"black_points": g.BlackPoints,
					// 		"winner":       g.WinnerIndex,
					// 		"reason":       g.GameEndedReason,
					// 		"moveNum":      g.MoveNum,
					// 	},
					// }
					g.emitGameEnded(g.MoveNum)

					g.mu.Lock()
					g.end()
					g.mu.Unlock()

					delete(PlayerToGame, g.Players[0])
					delete(PlayerToGame, g.Players[1])

					saveGameToDB(g)
					return
				}

			case <-g.MovePlayed:
				g.mu.Lock()
				turn = g.CurrectTurn
				ended := g.GameProgress == GAME_ENDED
				g.mu.Unlock()

				if ended {
					return
				}
				g.Clock.Switch(turn)
			case <-g.Done():
				return
			}
		}
	}()
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
		return fmt.Errorf("Not your turn")
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

	g.MoveNum++
	g.MovePlayed <- struct{}{}

	saveGameToDB(g)
	saveMoveToDB(g, m)
	saveSnapshotIfNeededToDB(g)

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

func (g *Game) Resign(player string) error {
	if (g.CurrectTurn == Black && player != g.Players[0]) ||
		(g.CurrectTurn == White && player != g.Players[1]) {
		return fmt.Errorf("Not your turn")
	}

	g.GameProgress = GAME_ENDED
	g.GameEndedReason = "resignation"
	if g.CurrectTurn == Black {
		g.WinnerIndex = 1
	}
	if g.CurrectTurn == White {
		g.WinnerIndex = 0
	}
	g.emitGameEnded(g.MoveNum)
	g.end()

	delete(PlayerToGame, g.Players[0])
	delete(PlayerToGame, g.Players[1])
	saveGameToDB(g)
	return nil
}

func (g *Game) endGame() {
	g.GameProgress = GAME_ENDED
	g.WhitePoints, g.BlackPoints = getPointsFromBoard(g.Board)
	g.GameEndedReason = "normal"
	if g.WhitePoints > g.BlackPoints {
		g.WinnerIndex = 1
	} else if g.BlackPoints > g.WhitePoints {
		g.WinnerIndex = 0
	}

	g.emitGameEnded(g.MoveNum + 1)
	g.end()
}

func (g *Game) emitGameEnded(moveNum int) {
	g.Events <- map[string]any{
		"type": "game_ended",
		"data": map[string]any{
			"players":      g.Players,
			"white_points": g.WhitePoints,
			"black_points": g.BlackPoints,
			"winner":       g.WinnerIndex,
			"reason":       g.GameEndedReason,
			"moveNum":      moveNum,
		},
	}
}

func (g *Game) hashMove(placedStone Stone, otherColorDeletedChains [][]Stone) {
	g.hash = g.zobrist.ToggleStone(g.hash, placedStone)
	for i := range otherColorDeletedChains {
		for _, stone := range otherColorDeletedChains[i] {
			g.hash = g.zobrist.ToggleStone(g.hash, stone)
		}
	}
}

func (g *Game) Done() <-chan struct{} { return g.done }

func (g *Game) end() {
	g.doneOnce.Do(func() { close(g.done) })
}

func (g *Game) Print() {
	g.Board.PrintBoard()
}
