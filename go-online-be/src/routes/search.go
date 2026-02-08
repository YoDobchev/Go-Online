package routes

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YoDobchev/Go-Online/src/database"
	"github.com/go-chi/chi/v5"
)

func SearchRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/games", searchGamesHandler)
	return r
}

type gameListItem struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status"`
	BoardSize int    `json:"boardSize"`
	Ranked    bool   `json:"ranked"`
	Players   int    `json:"players"`
	CreatedAt string `json:"createdAt,omitempty"`
}

type gameListResponse struct {
	Games      []gameListItem `json:"games"`
	Page       int            `json:"page"`
	TotalPages int            `json:"totalPages"`
}

func computeStatus(g database.Game) string {
	if g.GameProgress >= 2 {
		return "finished"
	}
	if g.PlayerBlack == nil || g.PlayerWhite == nil {
		return "open"
	}
	return "running"
}

func computePlayers(g database.Game) int {
	n := 0
	if g.PlayerBlack != nil {
		n++
	}
	if g.PlayerWhite != nil {
		n++
	}
	return n
}

func computeName(g database.Game) string {
	switch {
	case g.PlayerBlack != nil && g.PlayerWhite != nil:
		return fmt.Sprintf("%s vs %s", *g.PlayerBlack, *g.PlayerWhite)
	case g.PlayerBlack != nil:
		return fmt.Sprintf("%s vs ?", *g.PlayerBlack)
	case g.PlayerWhite != nil:
		return fmt.Sprintf("? vs %s", *g.PlayerWhite)
	default:
		return ""
	}
}

func searchGamesHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	rankedStr := strings.TrimSpace(r.URL.Query().Get("ranked"))
	sizeStr := strings.TrimSpace(r.URL.Query().Get("size"))

	page := 1
	if ps := strings.TrimSpace(r.URL.Query().Get("p")); ps != "" {
		if n, err := strconv.Atoi(ps); err == nil && n >= 1 {
			page = n
		}
	}

	const pageSize = 10

	db := database.DB.WithContext(r.Context()).Model(&database.Game{})

	if q != "" {
		like := "%" + q + "%"
		db = db.Where("(player_black ILIKE ? OR player_white ILIKE ?)", like, like)
	}

	if sizeStr != "" {
		if n, err := strconv.Atoi(sizeStr); err == nil && (n == 9 || n == 13 || n == 19) {
			db = db.Where("board_size = ?", n)
		}
	}

	if rankedStr != "" {
		if b, err := strconv.ParseBool(rankedStr); err == nil {
			db = db.Where("ranked = ?", b)
		}
	}

	switch status {
	case "finished":
		db = db.Where("game_progress >= ?", 2)
	case "open":
		db = db.Where("(game_progress < ? AND (player_black IS NULL OR player_white IS NULL))", 2)
	case "running":
		db = db.Where("(game_progress < ? AND (player_black IS NOT NULL AND player_white IS NOT NULL))", 2)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, "failed to count games")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	var rows []database.Game
	if err := db.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error; err != nil {
		writeJSON(w, http.StatusInternalServerError, "failed to load games")
		return
	}

	out := make([]gameListItem, 0, len(rows))
	for _, g := range rows {
		item := gameListItem{
			ID:        g.ID,
			Name:      computeName(g),
			Status:    computeStatus(g),
			BoardSize: g.BoardSize,
			Ranked:    g.Ranked,
			Players:   computePlayers(g),
		}
		if !g.CreatedAt.IsZero() {
			item.CreatedAt = g.CreatedAt.UTC().Format(time.RFC3339)
		}
		out = append(out, item)
	}

	writeJSON(w, http.StatusOK, gameListResponse{
		Games:      out,
		Page:       page,
		TotalPages: totalPages,
	})
}
