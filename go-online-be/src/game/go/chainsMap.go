package gogame

import "fmt"

type Stone struct {
	x     int
	y     int
	color uint8
}

type Point struct {
	x int
	y int
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

func (c *ChainsMap) AddStone(x, y int, color uint8) ([]Stone, [][]Stone) {
	current := Stone{x, y, color}
	neighbors := c.getNeighbors(x, y)
	sameColor, oppositeColor := groupByColor(neighbors, color)
	sameColorCapturedChain := []Stone{}
	otherColorCapturedChains := [][]Stone{}

	if len(sameColor) == 0 {
		//Create a new chain
		index := c.LastIndex + 1
		c.Map[x][y] = index
		c.Liberties[index] = len(neighbors) - len(oppositeColor)
		if c.Liberties[index] == 0 {
			sameColorCapturedChain = []Stone{current}
		}
		c.LastIndex++
	} else {
		//Add to existing chain and merge chains if necessary
		index := c.Map[sameColor[0].x][sameColor[0].y]
		c.Map[x][y] = index
		sameColorCapturedChain = c.updateChainDFS(current)
	}

	seen := make(map[int]struct{})
	for _, o := range oppositeColor {
		idx := c.Map[o.x][o.y]
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		chain := c.updateChainDFS(o)
		if len(chain) > 0 {
			otherColorCapturedChains = append(otherColorCapturedChains, chain)
		}
	}

	return sameColorCapturedChain, otherColorCapturedChains
}

func (c *ChainsMap) Snapshot() ChainsMap {
	mapCopy := make([][]int, len(c.Map))
	for i := range c.Map {
		mapCopy[i] = append([]int(nil), c.Map[i]...)
	}

	libsCopy := make(map[int]int, len(c.Liberties))
	for k, v := range c.Liberties {
		libsCopy[k] = v
	}

	return ChainsMap{
		Map:       mapCopy,
		Liberties: libsCopy,
		LastIndex: c.LastIndex,
	}
}

func (c *ChainsMap) Restore(snapshot ChainsMap) {
	c.Map = snapshot.Map
	c.Liberties = snapshot.Liberties
	c.LastIndex = snapshot.LastIndex
}

func (c *ChainsMap) updateChainDFS(start Stone) []Stone {
	visited := make(map[Point]struct{})
	liberties := make(map[Point]struct{})
	color := start.color
	chainIndex := c.Map[start.x][start.y]

	stack := []Point{{start.x, start.y}}
	chain := []Stone{}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, seen := visited[p]; seen {
			continue
		}
		visited[p] = struct{}{}
		c.Map[p.x][p.y] = chainIndex
		chain = append(chain, Stone{p.x, p.y, color})

		for _, n := range c.getNeighbors(p.x, p.y) {
			switch n.color {
			case Empty:
				liberties[Point{n.x, n.y}] = struct{}{}
			case color:
				stack = append(stack, Point{n.x, n.y})
			}
		}
	}
	c.Liberties[chainIndex] = len(liberties)

	if len(liberties) == 0 {
		return chain
	}
	return []Stone{}
}

func (c *ChainsMap) applyCapture(sameColorCapturedChain []Stone, otherColorCapturedChains [][]Stone) {
	for _, s := range sameColorCapturedChain {
		c.Map[s.x][s.y] = 0
	}
	for i := range otherColorCapturedChains {
		for _, s := range otherColorCapturedChains[i] {
			c.Map[s.x][s.y] = 0
		}
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
