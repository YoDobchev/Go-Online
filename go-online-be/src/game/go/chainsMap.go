package gogame

type ChainsMap struct {
	Map       [][]uint64 // [x][y][group]
	Count     uint64
	Liberties map[uint64]uint64 // [group] -> num of liberties
}

func NewChainsMap(boardSize int) *ChainsMap {
	c := new(ChainsMap)
	c.Map = make([][]uint64, boardSize)
	for i := range c.Map {
		c.Map[i] = make([]uint64, boardSize)
	}
	c.Liberties = make(map[uint64]uint64)
	return c
}

func (c *ChainsMap) AddStone(x, y int, color uint8) {

}
