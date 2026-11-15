package main

import (
	"fmt"
	"os"
	"math"
	"sync"
)

type point struct { x, y float64 }
type edge struct { a, b int }
type channel struct {
	outs [] chan point
	ins [] chan point
}

func near(coords []point, r float64) []edge {
	var edges []edge

	for i, ipct := range coords {
		for j := i + 1; j < len(coords); j++ {
			jpct := coords[j]
			x_diff := ipct.x - jpct.x
			y_diff := ipct.y - jpct.y
			if x_diff*x_diff + y_diff*y_diff <= r*r {
				edges = append(edges, edge {i, j})
				edges = append(edges, edge {j, i})
			}
		}
	}

	return edges
}

func check(coords []point, precision float64) bool {
	fst := coords[0]
	for _, pct := range coords {
		if ( math.Abs(fst.x - pct.x) > precision || math.Abs(fst.y - pct.y) > precision ) {
			return true
		}
	}

	return false
}

func make_channels(edges []edge, n int) []channel {
	channels := make([]channel, n)

	for _, e := range edges {
		ch := make(chan point, 1)
		channels[e.a].outs = append(channels[e.a].outs, ch)
		channels[e.b].ins = append(channels[e.b].ins, ch)
	}

	return channels
}

func robo(
	wg 			*sync.WaitGroup,
	myCoord *point,
	ins 		[]chan point,
	outs 		[]chan point,
	id 			int) {
	defer wg.Done()
	sum_x := float64(0)
	sum_y := float64(0)

	// Look: send(pos); receive(vecini)
	for _, out := range outs {
		out <- *myCoord
	}

	for _, in := range ins {
		pct := <- in
		sum_x += pct.x
		sum_y += pct.y
	}

	// Compute: new_position = (pos + sum(vecini)) / (1 + len(vecini))
	new_x := (myCoord.x + sum_x) / float64(1 + len(ins))
	new_y := (myCoord.y + sum_y) / float64(1 + len(ins))

	// Move: pos = new_position
	*myCoord = point { new_x, new_y }
}

func main() {
	// read from file
	coords := read_coords("coords.txt")
	if coords == nil {
		panic(1)
	}
	fmt.Println(len(coords), coords)

	// read done. now start
	var wg sync.WaitGroup
	r := float64(15)
	limit := 10000
	precision := 1e-4

	// while they are not stacked + range limit
	
	for iter := 0; check(coords, precision) && iter < limit; iter++ {
		edges := near(coords, r)
		channels := make_channels(edges, len(coords))

		wg.Add(len(coords))
		for i := range len(coords) {
			go robo(&wg, &coords[i], channels[i].ins, channels[i].outs , i)
		}
		
		wg.Wait()
	}

	for _, coord := range coords {
		fmt.Printf("{%.1f, %.1f} ", coord.x, coord.y)
	}
	fmt.Println()
}

func read_coords (filename string) []point {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var nr int
	_, err = fmt.Fscan(file, &nr)
	if err != nil {
		return nil
	}

	coords := make([]point, 0, nr)
	for range nr {
		var pct point

		_, err = fmt.Fscan(file, &pct.x, &pct.y)
		if err != nil {
			return nil
		}

		coords = append(coords, pct)
	}

	return coords
}
