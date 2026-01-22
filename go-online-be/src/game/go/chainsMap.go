package gogame

import "fmt"

type Stone struct {
	x     int
	y     int
	color uint8
}

const MAX_NEIGHBORS = 4

type ChainsMap struct {
	Map       [][]int //[x][y] -> chain index starting from 1 (0 if empty)
	LastIndex int
	Liberties map[int]int //[chain index] -> num of liberties

	board *Board
}

func NewChainsMap(board *Board) *ChainsMap {
	c := new(ChainsMap)
	c.Map = make([][]int, board.Size)
	for i := range c.Map {
		c.Map[i] = make([]int, board.Size)
	}
	c.Liberties = make(map[int]int)
	c.board = board
	return c
}

func (c *ChainsMap) AddStone(x, y int, color uint8) {
	neighbors := c.getNeighbors(x, y)
	sameColor, oppositeColor := groupByColor(neighbors, color)

	liberties := len(neighbors) - len(oppositeColor) - len(sameColor)

	//fmt.Printf("Same color %v", len(sameColor))
	fmt.Println()
	if len(sameColor) == 0 {
		//Create a new chain
		index := c.LastIndex + 1
		c.Map[x][y] = index

		//Calculate liberties
		c.Liberties[index] = liberties
		c.LastIndex++
	} else if len(sameColor) == 1 {
		//Add to existing chain
		neighbor := sameColor[0]
		index := c.Map[neighbor.x][neighbor.y]
		c.Map[x][y] = index

		//Update liberties
		c.Liberties[index] += (liberties - 1)
	} else {
		c.unifyChains(sameColor, liberties)
	}

}

func groupByColor(neighbors []Stone, color uint8) ([]Stone, []Stone) {
	sameColor := []Stone{}
	oppositeColor := []Stone{}
	for _, n := range neighbors {
		if n.color == Empty {
			continue
		}

		if n.color == color {
			sameColor = append(sameColor, n)
		} else {
			oppositeColor = append(oppositeColor, n)
		}
	}
	return sameColor, oppositeColor
}

func (c *ChainsMap) unifyChains(sameColor []Stone, liberties int) {
	neighborGroups := make(map[int]struct{})
	for _, n := range sameColor {
		group := c.Map[n.x][n.y]
		_, found := neighborGroups[group]
		if !found {
			neighborGroups[group] = struct{}{}
		}
	}

	unifiedGroup := c.Map[sameColor[0].x][sameColor[0].y]

	//Calculate liberties
	libertiesSum := liberties - len(neighborGroups)
	for i := range neighborGroups {
		libertiesSum += c.Liberties[i]
	}
	c.Liberties[unifiedGroup] = libertiesSum

	//Reassign group indexes to new unifiedGroup index
	for i := range c.Map {
		for j := range c.Map[i] {
			oldGroup := c.Map[i][j]
			_, found := neighborGroups[oldGroup]
			if found {
				c.Map[i][j] = unifiedGroup
			}
		}
	}
}

func (c *ChainsMap) getNeighbors(x, y int) []Stone {
	neighbors := []Stone{}
	neighbors = c.appendStone(x-1, y, neighbors)
	neighbors = c.appendStone(x+1, y, neighbors)
	neighbors = c.appendStone(x, y-1, neighbors)
	neighbors = c.appendStone(x, y+1, neighbors)
	fmt.Println(len(neighbors))
	return neighbors
}

func (c *ChainsMap) appendStone(x, y int, stones []Stone) []Stone {
	square, err := c.board.GetSquare(x, y)
	if err == nil {
		stones = append(stones, Stone{x, y, square.stone})
	}
	return stones
}

func (c *ChainsMap) PrintGroups() {
	for i := range c.Map {
		for j := range c.Map[i] {
			fmt.Printf("%v ", c.Map[i][j])
		}
		fmt.Println()
	}
}
