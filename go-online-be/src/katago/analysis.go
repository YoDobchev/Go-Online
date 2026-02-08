package katago

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Engine struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
}

type QueryMove struct {
	Player string `json:"player"`
	Move   string `json:"move"`
}

type Query struct {
	ID         string      `json:"id"`
	BoardXSize int         `json:"boardXSize"`
	BoardYSize int         `json:"boardYSize"`
	Rules      string      `json:"rules"`
	Komi       float32     `json:"komi"`
	Moves      [][2]string `json:"moves"`
	MaxVisits  int         `json:"maxVisits"`
}

type Response struct {
	ID       string `json:"id"`
	RootInfo struct {
		Winrate       float64 `json:"winrate"`
		CurrentPlayer string  `json:"currentPlayer"`
	} `json:"rootInfo"`
}

var Eng *Engine

func InitEng() {
	bin := os.Getenv("KATAGO_BIN")
	model := os.Getenv("KATAGO_ANALYSIS_MODEL")
	cfg := os.Getenv("KATAGO_ANALYSIS_CONFIG")
	fmt.Println(bin + " " + model + " " + cfg)
	if bin == "" || model == "" || cfg == "" {
		panic("katago env not configured")
	}
	var err error
	Eng, err = Start(context.Background(), bin, cfg, model)
	if err != nil {
		panic(fmt.Sprintf("failed to start KataGo engine: %v", err))
	}
}

func Start(ctx context.Context, katagoPath, cfgPath, modelPath string) (*Engine, error) {
	cmd := exec.CommandContext(ctx, katagoPath,
		"analysis",
		"-config", cfgPath,
		"-model", modelPath,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(stdoutPipe)
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 10*1024*1024)

	return &Engine{
		cmd:    cmd,
		stdin:  stdin,
		stdout: sc,
	}, nil
}

func (e *Engine) Analyze(ctx context.Context, q Query) (Response, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	b, err := json.Marshal(q)
	if err != nil {
		return Response{}, err
	}
	if _, err := e.stdin.Write(append(b, '\n')); err != nil {
		return Response{}, err
	}

	for {
		select {
		case <-ctx.Done():
			return Response{}, ctx.Err()
		default:
		}
		if !e.stdout.Scan() {
			if err := e.stdout.Err(); err != nil {
				return Response{}, err
			}
			return Response{}, errors.New("katago: stdout closed")
		}
		var resp Response
		if err := json.Unmarshal(e.stdout.Bytes(), &resp); err != nil {
			continue
		}
		if resp.ID == q.ID {
			return resp, nil
		}
	}
}

func SGFCoord(x, y int) string {
	return string([]byte{byte('a' + x), byte('a' + y)})
}

func BW(color uint8) string {
	if color == 2 {
		return "b"
	}
	return "w"
}

func BlackWinProb(resp Response) float64 {
	bw := resp.RootInfo.Winrate
	if resp.RootInfo.CurrentPlayer == "W" {
		return 1.0 - bw
	}
	return bw
}

func GTPCoord(x, y, size int) string {
	const letters = "ABCDEFGHJKLMNOPQRST"
	if x < 0 || y < 0 || x >= size || y >= size {
		return "pass"
	}
	return fmt.Sprintf("%c%d", letters[y], size-x)
}
