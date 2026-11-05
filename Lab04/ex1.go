package main
import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
)

type node struct {
	in <-chan int
	out chan<- int
	index int
	leader int
	isLider bool
	wg *sync.WaitGroup
}

func (n node) start() {
	defer n.wg.Done()

	n.out <- n.index
	for recv := range n.in {
		if recv < n.leader {
			continue
		} else if recv == n.index {
			fmt.Println(n.index, " leader!")
			n.isLider = true
			break
		} else {

			n.leader = recv
			n.out <- n.leader
			fmt.Println(n.index,": recv:", recv, "lead:", n.leader)
		}
	}

	close(n.out)
}

func main() {
	var wg sync.WaitGroup
	n := 10
	nodes := make([]node, 0, n)
	ids_used := make([]int, 0, n)

	in := make(chan int, 1)
	initial := in
	out := make(chan int, 1)

	// make nodes
	wg.Add(n)
	for range n {
		id := make_id(&ids_used)
		crr_node := node{in, out, id, id, false, &wg}
		nodes = append(nodes, crr_node)
		in = out
		out = make(chan int, 1)
	}
	nodes[len(nodes) - 1].out = initial
	fmt.Println("Ids", len(ids_used), ":", ids_used)

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
