package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type PlayMovePayload struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type client struct {
	conn *websocket.Conn
	role string
	user string
}

type gameHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]*client
}

func newGameHub() *gameHub {
	return &gameHub{clients: make(map[*websocket.Conn]*client)}
}

func (h *gameHub) add(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.conn] = c
}

func (h *gameHub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *gameHub) broadcast(v any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		_ = conn.WriteJSON(v)
	}
}

var (
	hubsMu sync.Mutex
	hubs   = map[string]*gameHub{}
)

func getHub(gameID string) *gameHub {
	hubsMu.Lock()
	defer hubsMu.Unlock()
	h, ok := hubs[gameID]
	if !ok {
		h = newGameHub()
		hubs[gameID] = h
	}
	return h
}

func WsGameHandler(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "id")
	game, exists := gogame.GameInstances[gameID]
	if !exists {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	if game.Players[0] == "" || game.Players[1] == "" {
		http.Error(w, "game not started yet", http.StatusForbidden)
		return
	}

	username := ""
	if user, err := middleware.GetUserInfo(r); err == nil {
		username = user.Username
	}

	role := "spectator"
	if username != "" && (username == game.Players[0] || username == game.Players[1]) {
		role = "player"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	hub := getHub(gameID)
	hub.add(&client{conn: conn, role: role, user: username})
	defer hub.remove(conn)

	_ = conn.WriteJSON(map[string]any{
		"type": "hello",
		"data": map[string]any{
			"role": role,
		},
	})

	_ = conn.WriteJSON(map[string]any{
		"type": "game_snapshot",
		"data": map[string]any{
			"players": game.Players,
			"turn":    game.CurrectTurn,
			"board":   game.Board,
		},
	})

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "play_move":
			if role != "player" {
				_ = conn.WriteJSON(map[string]any{
					"type": "error",
					"data": "spectators cannot play moves",
				})
				continue
			}

			var p PlayMovePayload
			if err := json.Unmarshal(msg.Data, &p); err != nil {
				_ = conn.WriteJSON(map[string]any{"type": "error", "data": "bad payload"})
				continue
			}

			err := game.PlayMove(username, p.Row, p.Col)
			if err != nil {
				_ = conn.WriteJSON(map[string]any{"type": "error", "data": err.Error()})
				continue
			}

			// hub.broadcast(map[string]any{
			// 	"type": "move_played",
			// 	"data": map[string]any{
			// 		"by":  username,
			// 		"row": p.Row,
			// 		"col": p.Col,
			// 	},
			// })
			hub.broadcast(map[string]any{
				"type": "game_snapshot",
				"data": map[string]any{
					"players": game.Players,
					"turn":    game.CurrectTurn,
					"board":   game.Board,
				},
			})
		default:
			_ = conn.WriteJSON(map[string]any{"type": "error", "data": "unknown message type"})
		}

		if game.GameProgress == gogame.GAME_ENDED {
			hub.broadcast(map[string]any{
				"type": "game_ended",
				"data": map[string]any{
					"white_points": game.WhitePoints,
					"black_points": game.BlackPoints,
				},
			})

			hub.mu.Lock()
			for conn := range hub.clients {
				conn.Close()
			}
			hub.clients = make(map[*websocket.Conn]*client)
			hub.mu.Unlock()
			break
		}
	}
}
