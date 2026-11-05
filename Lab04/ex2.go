package main
import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
)

type node struct {
	in []<-chan int
	out []chan<- int
	id int
	max_id int
	isLider bool
	rounds int
	diam int
	wg *sync.WaitGroup
}

func (n node) start() {
	defer n.wg.Done()

	for n.rounds != n.diam {
		n.rounds ++

		for _, ch := range n.out {
			ch <- n.max_id
		}

		for _, ch := range n.in {
			n.max_id = max(n.max_id, <-ch)
		}
	}

	if n.max_id == n.id {
		fmt.Println(n.id, " leader!")
		n.isLider = true
	}
}

func main() {
	var wg sync.WaitGroup
	n := 5
	diam := 3

	nodes := make([]node, n)
	ids_used := make([]int, 0, n)

	// init nodes 
	for i := range n {
		id := make_id(&ids_used)
		nodes[i] = node{
			in:      make([]<-chan int, 0),
			out:     make([]chan<- int, 0),
			id:      id,
			max_id:  id,
			isLider: false,
			rounds:  0,
			diam:    diam,
			wg:      &wg,
		}
	}
	fmt.Println("Ids", len(ids_used), ":", ids_used)

	// make connections
	edges := [][2]int{
		{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 0},
		{0, 2}, {1, 3}, {2, 4}, {3, 0}, {4, 1},
	}

	// match channels
	for _, e := range edges {
		u, v := e[0], e[1]
		ch := make(chan int, 1)
		nodes[u].out = append(nodes[u].out, ch)
		nodes[v].in = append(nodes[v].in, ch)
	}

	// start
	wg.Add(n)
	for _, crr_node := range nodes {
		go crr_node.start()
	}

	wg.Wait()
	fmt.Println("Done")
}

func make_id (ids_used *[]int) int {
	maxim := 100
	idx := rand.Intn(maxim)
	
	for slices.Contains(*ids_used, idx) {
		idx = rand.Intn(maxim)
	}
	(*ids_used) = append((*ids_used), idx)

	return idx
}
