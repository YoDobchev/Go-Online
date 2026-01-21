package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
)

var (
	quit bool = false
	g    *gogame.Game
)

func parseCmd(cmd string) error {
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return fmt.Errorf("%s", "empty cmd")
	}
	switch parts[0] {
	case "quit":
		quit = true
	case "play":
		if len(parts) != 3 {
			return fmt.Errorf("%s", "too many of few move args")
		}

		err := execMove([2]string{parts[1], parts[2]})

		if err != nil {
			return err
		}

	case "pass":
		execMove([2]string{"-1", "-1"})
	default:
		return fmt.Errorf("unrecognized cmd")
	}

	return nil
}

func execMove(coords [2]string) error {

	x, err := strconv.Atoi(coords[0])
	if err != nil {
		return err
	}

	y, err := strconv.Atoi(coords[1])
	if err != nil {
		return err
	}

	err = g.PlayMove(x, y)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var err error
	g, err = gogame.NewGame(9, [2]string{"black", "white"})
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for !quit {
		g.PrintBoard()

		var playerTurnName string
		if g.CurrectTurn == gogame.Black {
			playerTurnName = g.Players[0]
		} else {
			playerTurnName = g.Players[1]
		}
		fmt.Printf("Enter cmd [%s]: ", playerTurnName)
		cmd, _ := reader.ReadString('\n')
		err := parseCmd(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}
