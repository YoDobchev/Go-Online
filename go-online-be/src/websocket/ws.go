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
	mu       sync.Mutex
	clients  map[*websocket.Conn]*client
	pumpOnce sync.Once
}

func (h *gameHub) startPump(game *gogame.Game) {
	h.pumpOnce.Do(func() {
		go func() {
			for {
				select {
				case ev, ok := <-game.Events:
					if !ok {
						return
					}
					h.broadcast(ev)
				case <-game.Done():
					return
				}
			}
		}()
	})
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

	if game.GameProgress == gogame.GAME_ENDED {
		http.Error(w, "game ended", http.StatusGone)
		return
	}

	username := ""
	if user, err := middleware.GetUserInfo(r); err == nil {
		username = user.Username
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	hub := getHub(gameID)
	c := &client{conn: conn, role: "spectator", user: username}
	hub.add(c)
	defer hub.remove(conn)

	role := "spectator"
	seat := "spectator"

	if username != "" {
		switch username {
		case game.Players[0]:
			role = "player"
			seat = "black"
		case game.Players[1]:
			role = "player"
			seat = "white"
		default:
			if game.GameProgress == gogame.GAME_WAITING_FOR_PLAYER &&
				(game.Players[0] == "" || game.Players[1] == "") {

				if err := game.Join(username); err == nil {
					switch username {
					case game.Players[0]:
						role = "player"
						seat = "black"
					case game.Players[1]:
						role = "player"
						seat = "white"
					}

					hub.broadcast(map[string]any{
						"type": "sync",
						"data": map[string]any{
							"players": game.Players,
							"turn":    game.CurrectTurn,
							"board":   game.Board,
							"moveNum": game.MoveNum,
						},
					})
				}
			}
		}
	}

	c.role = role

	_ = conn.WriteJSON(map[string]any{
		"type": "hello",
		"data": map[string]any{
			"role": role,
			"seat": seat, // "black"|"white"|"spectator"
		},
	})

	_ = conn.WriteJSON(map[string]any{
		"type": "sync",
		"data": map[string]any{
			"players": game.Players,
			"turn":    game.CurrectTurn,
			"board":   game.Board,
			"moveNum": game.MoveNum,
		},
	})

	hub.startPump(game)

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

			if err := game.PlayMove(username, p.Row, p.Col); err != nil {
				_ = conn.WriteJSON(map[string]any{"type": "error", "data": err.Error()})
				continue
			}

			hub.broadcast(map[string]any{
				"type": "sync",
				"data": map[string]any{
					"players": game.Players,
					"turn":    game.CurrectTurn,
					"board":   game.Board,
					"moveNum": game.MoveNum,
				},
			})
		default:
			_ = conn.WriteJSON(map[string]any{"type": "error", "data": "unknown message type"})
		}

		if game.GameProgress == gogame.GAME_ENDED {
			if !game.EndedByTimeout {
				hub.broadcast(map[string]any{
					"type": "game_ended",
					"data": map[string]any{
						"white_points": game.WhitePoints,
						"black_points": game.BlackPoints,
						"winner":       game.WinnerIndex,
					},
				})
			}

			delete(gogame.PlayerToGame, game.Players[0])
			delete(gogame.PlayerToGame, game.Players[1])

			hub.mu.Lock()
			for cconn := range hub.clients {
				_ = cconn.Close()
			}
			hub.clients = make(map[*websocket.Conn]*client)
			hub.mu.Unlock()

			break
		}
	}
}
