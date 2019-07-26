package main

import "fmt"
import "os"
import "bufio"
import "bytes"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type node struct {
	x      int
	y      int
	branch bool
}

type maze struct {
	layers    [][]byte
	depth     int
	entrance  node
	nodes     []node
	exit      node
	solutions [][]node
	best      []node
}

type path struct {
	current    node
	lastBranch node
	behind     node
	allTravel  []node
	clearRooms [4]node
	bestTime   []node
	branches   []node
}

type compass struct {
	N bool
	S bool
	E bool
	W bool
}

func bFile(f string) []byte {
	mazeFile := f
	file, err := os.Open(mazeFile)
	check(err)
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bMap := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bMap)
	check(err)
	return bMap
}

func (m *maze) findEntrance() {
	for y := range m.layers {
		if m.layers[y][0] == 32 {
			m.entrance = node{0, y, false}
			return
		}
	}
	fmt.Println("No Entrance!")
	return
}

func (m *maze) findExit() {
	y := len(m.layers) - 2
	for x := range m.layers[(0)] {
		if m.layers[y][(x)] == 32 {
			m.exit = node{x, y, true}
			return
		}
	}
	fmt.Println("No Entrance!")
	return
}

func (p *path) getSurronding(m maze) {
	ymax := len(m.layers) - 2
	xmax := len(m.layers[0])

	north := node{x: p.current.x, y: (p.current.y - 1)}
	south := node{x: p.current.x, y: (p.current.y + 1)}
	east := node{x: (p.current.x + 1), y: (p.current.y)}
	west := node{x: (p.current.x - 1), y: (p.current.y)}
	nsew := [4]node{north, south, east, west}
	var total int
	for r := range nsew {

		if nsew[r].x >= 0 && nsew[r].x <= xmax && nsew[r].y >= 0 && nsew[r].y <= ymax {

			if m.layers[nsew[r].y][nsew[r].x] == 32 {
				total += 1
				p.clearRooms[r] = node{nsew[r].x, nsew[r].y, true}
				p.branches = append(p.branches, node{nsew[r].x, nsew[r].y, true})
				p.lastBranch = p.clearRooms[r]
			} else {
				p.clearRooms[r] = node{nsew[r].x, nsew[r].y, false}
			}
		}
	}
	if total > 2 {
		p.current.branch = true
		return
	} else {
		return
	}

}
func (p *path) moveForward(m maze) {
	if p.current == m.entrance {
		p.allTravel = append(p.allTravel, p.current)
	}
	if p.current == m.exit {
		p.endGame(m)
	}
	p.getSurronding(m)
	if p.current.branch != true {
		for i := range p.clearRooms {

			if (p.clearRooms[i].branch) && (node{p.behind.x, p.behind.y, false} != node{p.clearRooms[i].x, p.clearRooms[i].y, false}) {
				p.behind = p.current
				p.current = node{p.clearRooms[i].x, p.clearRooms[i].y, false}
				p.allTravel = append(p.allTravel, p.current)

				fmt.Println("Appending: ", p.current)
				if p.current.x == m.exit.x && p.current.y == m.exit.y {
					p.endGame(m)
				}
				break
			}

		}
	} else {
		for i := range p.clearRooms {
			if p.clearRooms[i].branch && (node{p.behind.x, p.behind.y, false}) != (node{p.clearRooms[i].x, p.clearRooms[i].y, false}) {
				p.behind = p.current
				p.current = node{p.clearRooms[i].x, p.clearRooms[i].y, false}
				p.allTravel = append(p.allTravel, p.current)
				fmt.Println("Appending: ", p.current)
				if p.current.x == m.exit.x && p.current.y == m.exit.y {
					p.endGame(m)
				}

			}
		}
	}
	p.moveForward(m)
}
func (m maze) findAllBranches() {
	// TODO: findAllBranches abandons the pretext of a maze goer, looks at the overhead view of the map
	//then determins all potentiall branches

}
func (p *path) smartPath() {
	//TODO: The smart path function will try to sole the maxe with some reasoning
	//keeping track of branches and, realizing when repeating steps, And recognizing
	//when it is important to backtrack
}

func (p *path) shortestDistance() {
	//TODO: Using findallBranches, finds the shortest distance to the end
	//Utilizes findAllBranches
}

func (m *maze) drawPath(p *path) [][]byte {
	finishedMap := m.layers
	for i := range p.bestTime {
		finishedMap[p.bestTime[i].y][p.bestTime[i].x] = byte(43) //[]byte("+")
	}
	return finishedMap
}
func (p *path) startPath(m maze) {
	p.current = m.entrance
	p.moveForward(m)

}
func createMaze(bm []byte) (m maze) {
	m = maze{
		layers: bytes.SplitAfterN(bm, []byte("\n"), -1),
		depth:  len(m.layers),
	}
	m.findEntrance()
	m.findAllBranches()
	m.findExit()
	return
}
func (p *path) endGame(m maze) {
	fmt.Println("GAMEOVER MAN")
	fmt.Println("You Won")
	p.bestTime = p.allTravel
	fMap := m.drawPath(p)
	for i := range fMap {
		fmt.Println(string(fMap[i]))
	}
	os.Exit(0)
}
func main() {
	bMap := bFile(os.Args[1])
	Maze := createMaze(bMap)
	newPath := path{}
	newPath.startPath(Maze)
	fmt.Println(newPath)

}
