package gogame

import (
	"fmt"
	"os"
	"sync"

	"github.com/YoDobchev/Go-Online/src/ai/gtpcoords"
	"github.com/YoDobchev/Go-Online/src/ai/katago"
	"github.com/YoDobchev/Go-Online/src/database"
)

const AIUsername = "KataGo"

var (
	aiMu       sync.Mutex
	aiEngines  = map[string]*katago.Engine{}
	aiThinking = map[string]bool{}
)

func gtpColor(turn uint8) string {
	if turn == Black {
		return "B"
	}
	return "W"
}

func isAITurn(g *Game) bool {
	return (g.CurrectTurn == Black && g.Players[0] == AIUsername) ||
		(g.CurrectTurn == White && g.Players[1] == AIUsername)
}

func ensureEngine(g *Game) (*katago.Engine, error) {
	aiMu.Lock()
	if e, ok := aiEngines[g.ID]; ok {
		aiMu.Unlock()
		return e, nil
	}
	aiMu.Unlock()

	bin := os.Getenv("KATAGO_BIN")
	model := os.Getenv("KATAGO_MODEL")
	cfg := os.Getenv("KATAGO_CONFIG")
	if bin == "" || model == "" || cfg == "" {
		return nil, fmt.Errorf("missing KATAGO_BIN/KATAGO_MODEL/KATAGO_CONFIG")
	}

	e, err := katago.Start(bin, model, cfg)
	if err != nil {
		return nil, err
	}

	if err := e.Init(g.Board.Size, g.Komi); err != nil {
		_ = e.Close()
		return nil, err
	}

	if database.DB != nil {
		var moves []database.GameMove
		_ = database.DB.Where("game_id = ?", g.ID).Order("move_no ASC").Find(&moves).Error
		for _, mv := range moves {
			col := "W"
			if mv.Color == Black {
				col = "B"
			}
			coord := "pass"
			if mv.X != PASS {
				c, err := gtpcoords.ToGTP(g.Board.Size, mv.X, mv.Y)
				if err != nil {
					_ = e.Close()
					return nil, err
				}
				coord = c
			}
			if err := e.Play(col, coord); err != nil {
				_ = e.Close()
				return nil, err
			}
		}
	}

	aiMu.Lock()
	aiEngines[g.ID] = e
	aiMu.Unlock()
	return e, nil
}

func StopEngine(gameID string) {
	aiMu.Lock()
	e := aiEngines[gameID]
	delete(aiEngines, gameID)
	delete(aiThinking, gameID)
	aiMu.Unlock()
	if e != nil {
		_ = e.Close()
	}
}

func tryLockThinking(gameID string) bool {
	aiMu.Lock()
	defer aiMu.Unlock()
	if aiThinking[gameID] {
		return false
	}
	aiThinking[gameID] = true
	return true
}

func unlockThinking(gameID string) {
	aiMu.Lock()
	delete(aiThinking, gameID)
	aiMu.Unlock()
}

func MaybePlayAIMove(g *Game) (played bool, err error) {
	if !g.VsAI || g.GameProgress != GAME_IN_PROGRESS {
		return false, nil
	}
	if !isAITurn(g) {
		return false, nil
	}
	if !tryLockThinking(g.ID) {
		return false, nil
	}
	defer unlockThinking(g.ID)

	e, err := ensureEngine(g)
	if err != nil {
		return false, err
	}

	color := gtpColor(g.CurrectTurn)
	coord, err := e.GenMove(color)
	if err != nil {
		return false, err
	}

	if coord == "pass" || coord == "PASS" {
		err = g.PlayMove(AIUsername, PASS, PASS)
		return err == nil, err
	}

	row, col, pass, err := gtpcoords.FromGTP(g.Board.Size, coord)
	if err != nil {
		return false, err
	}
	if pass {
		err = g.PlayMove(AIUsername, PASS, PASS)
		return err == nil, err
	}

	err = g.PlayMove(AIUsername, row, col)
	return err == nil, err
}
