package katago

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Engine struct {
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
	mu  sync.Mutex
}

func Start(bin, model, cfg string) (*Engine, error) {
	cmd := exec.Command(bin, "gtp", "-model", model, "-config", cfg)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Engine{cmd: cmd, in: stdin, out: bufio.NewReader(stdout)}, nil
}

func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_ = e.in.Close()
	if e.cmd.Process != nil {
		return e.cmd.Process.Kill()
	}
	return nil
}

func (e *Engine) send(cmd string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := io.WriteString(e.in, cmd+"\n"); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for {
		line, err := e.out.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	resp := strings.TrimSpace(buf.String())
	if strings.HasPrefix(resp, "?") {
		return "", fmt.Errorf("gtp error: %s", resp)
	}
	if !strings.HasPrefix(resp, "=") {
		return "", fmt.Errorf("unexpected gtp response: %s", resp)
	}
	return strings.TrimSpace(strings.TrimPrefix(resp, "=")), nil
}

func (e *Engine) Init(size int, komi float32) error {
	if _, err := e.send("clear_board"); err != nil {
		return err
	}
	if _, err := e.send(fmt.Sprintf("boardsize %d", size)); err != nil {
		return err
	}
	if _, err := e.send(fmt.Sprintf("komi %.1f", komi)); err != nil {
		return err
	}
	return nil
}

func (e *Engine) Play(color, coord string) error {
	_, err := e.send(fmt.Sprintf("play %s %s", color, coord))
	return err
}

func (e *Engine) GenMove(color string) (string, error) {
	return e.send(fmt.Sprintf("genmove %s", color))
}
