package gogame

func getPointsFromBoard(board *Board) (int, int) {
	regionMap := getRegionMap(board)
	var white, black int
	for i := range regionMap {
		for j := range regionMap[i] {
			switch regionMap[i][j].stone {
			case White:
				white++
			case Black:
				black++
			}
		}
	}
	return white, black
}

func getRegionMap(board *Board) [][]Square {
	regionMap := initMap(board.Size)
	visited := make(map[Point]struct{})

	for i := range board.Squares {
		for j := range board.Squares[i] {
			color := board.Squares[i][j].stone
			if color != 0 {
				regionMap[i][j] = Square{color}
			} else if _, seen := visited[Point{i, j}]; !seen {
				visited, regionMap = colorRegions(i, j, board, visited, regionMap)
			}
		}
	}
	return regionMap
}

func colorRegions(x int, y int, board *Board, visited map[Point]struct{}, regionMap [][]Square) (map[Point]struct{}, [][]Square) {
	queue := []Point{{x, y}}
	region := []Point{{x, y}}
	seenWhite := false
	seenBlack := false

	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]

		if _, seen := visited[current]; seen {
			continue
		}

		for _, n := range board.getNeighbors(current.x, current.y) {
			switch n.color {
			case Empty:
				point := Point{n.x, n.y}
				queue = append(queue, point)
				region = append(region, point)
			case White:
				seenWhite = true
			case Black:
				seenBlack = true
			}
		}
		visited[current] = struct{}{}
	}

	var fillColor uint8
	if seenWhite && seenBlack {
		fillColor = Empty
	} else if seenWhite {
		fillColor = White
	} else if seenBlack {
		fillColor = Black
	} else {
		fillColor = Empty
	}

	for _, p := range region {
		regionMap[p.x][p.y] = Square{fillColor}
	}
	return visited, regionMap
}

func initMap(size int) [][]Square {
	regionMap := make([][]Square, size)
	for i := range regionMap {
		regionMap[i] = make([]Square, size)
	}
	return regionMap
}
