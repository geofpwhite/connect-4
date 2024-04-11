package connect4

import "fmt"

const EMPTY = 0
const BLUE = 1
const RED = 2

var Game *connect4 = newGame()

// field[0][0] is the bottom left of the board
type connect4 struct {
	field [][]int
	turn  int
}

func newGame() *connect4 {
	c4 := &connect4{
		field: make([][]int, 8),
		turn:  BLUE,
	}
	for i := 0; i < 8; i++ {
		c4.field[i] = make([]int, 8)
	}
	return c4
}

func (c4 *connect4) Clear() {
	c4.field = make([][]int, 8)
	c4.turn = BLUE
	for i := 0; i < 8; i++ {
		c4.field[i] = make([]int, 8)
	}
}

func (c4 *connect4) Insert(team, row int) bool {
	defer c4.Fall()
	defer println(c4.field[7][row])
	if c4.field[7][row] == 0 {
		c4.field[7][row] = team
		c4.turn = (team % 2) + 1
		return true
	}
	return false
}

func (c4 *connect4) Rotate() {
	defer c4.Fall()

	size := len(c4.field)
	newField := make([][]int, size)

	for i := range newField {
		for j := range newField {
			newField[size-j-1] = append(newField[size-j-1], c4.field[i][j])
		}
	}
	c4.field = newField
}

func (c4 *connect4) Fall() {
	for j := range c4.field {
		for i := range c4.field {
			k := i
			for k > 0 && c4.field[k-1][j] == EMPTY && c4.field[k][j] != EMPTY {
				c4.field[k-1][j] = c4.field[k][j]
				c4.field[k][j] = EMPTY
				k--
			}
		}
	}
}

type queueElement struct {
	team            int
	coordinate      [2]int
	upStreak        int
	rightStreak     int
	rightUpStreak   int
	rightDownStreak int
}

// rotating may cause both players to have
func (c4 *connect4) scanForConnect4() (int, map[queueElement]bool) {
	winner := EMPTY
	winners := map[queueElement]bool{}
	defer println(winner)
	defer fmt.Println(winners)
	coordinateQueue := []queueElement{}
	for i, num := range c4.field[0] {
		if num > 0 {
			coordinateQueue = append(coordinateQueue, queueElement{team: num, coordinate: [2]int{0, i}})
		}
	}
	// start  by checking bottom left , then check each neighbor going to the right, going up, going down, but not going left.

	for len(coordinateQueue) > 0 {
		poppedElement := coordinateQueue[0]
		coordinateQueue = coordinateQueue[1:]
		if poppedElement.upStreak >= 3 || poppedElement.rightStreak >= 3 || poppedElement.rightUpStreak >= 3 || poppedElement.rightDownStreak >= 3 {
			winners[poppedElement] = true
		}
		//check 4 neighbors
		up := [2]int{poppedElement.coordinate[0] + 1, poppedElement.coordinate[1]}
		right := [2]int{poppedElement.coordinate[0], poppedElement.coordinate[1] + 1}
		rightUp := [2]int{poppedElement.coordinate[0] + 1, poppedElement.coordinate[1] + 1}
		rightDown := [2]int{poppedElement.coordinate[0] - 1, poppedElement.coordinate[1] + 1}
		possibleNeighborCoordinates := [][2]int{up, right, rightUp, rightDown}
		for i, coord := range possibleNeighborCoordinates {
			if coord[0] >= 0 && coord[0] < 8 && coord[1] >= 0 && coord[1] < 8 {
				neighbor := queueElement{coordinate: coord, team: c4.field[coord[0]][coord[1]]}
				if neighbor.team == poppedElement.team {
					if i == 0 {
						neighbor.upStreak = poppedElement.upStreak + 1
						neighbor.rightStreak = 0
						neighbor.rightUpStreak = 0
						neighbor.rightDownStreak = 0
					} else if i == 1 {
						neighbor.rightStreak = poppedElement.rightStreak + 1
						neighbor.rightDownStreak = 0
						neighbor.rightUpStreak = 0
						neighbor.upStreak = 0
					} else if i == 2 {
						neighbor.rightUpStreak = poppedElement.rightUpStreak + 1
						neighbor.rightDownStreak = 0
						neighbor.rightStreak = 0
						neighbor.upStreak = 0
					} else if i == 3 {
						neighbor.rightDownStreak = poppedElement.rightDownStreak + 1
						neighbor.rightUpStreak = 0
						neighbor.rightStreak = 0
						neighbor.upStreak = 0
					}

					coordinateQueue = append(coordinateQueue, neighbor)
				} else if neighbor.team != EMPTY {
					coordinateQueue = append(coordinateQueue, neighbor)
				}
			}
		}
	}
	return winner, winners
}

func Main() {
	serve()
}
