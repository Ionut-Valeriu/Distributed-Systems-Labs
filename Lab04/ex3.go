package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type node struct {
	id     int
	in     []<-chan float64
	out    []chan<- float64
	val    float64
	rounds int
	wg     *sync.WaitGroup
}

const maxRounds = 30

func (n *node) start() {
	defer n.wg.Done()

	for r := 0; r < n.rounds; r++ {
		// Send
		for _, ch := range n.out {
			ch <- n.val
		}

		// Recv
		sum := n.val
		for _, ch := range n.in {
			sum += <-ch
		}

		// Local update
		deg := float64(len(n.in))
		n.val = sum / (deg + 1.0)
	}

	// Decision
	decision := 0
	if n.val > 0.5 {
		decision = 1
	}
	fmt.Printf("Node %d: final average = %.4f -> decision = %d\n", n.id, n.val, decision)
}

func main() {
	var wg sync.WaitGroup
	n := 5

	// Initialize nodes
	nodes := make([]*node, n)
	sumVotes := 0.0
	for i := range n {
		vote := float64(rand.Intn(2))
		sumVotes += vote
		nodes[i] = &node{
			id:     i,
			in:     make([]<-chan float64, 0),
			out:    make([]chan<- float64, 0),
			val:    vote,
			rounds: maxRounds,
			wg:     &wg,
		}
		fmt.Printf("Node %d initial vote: %.0f\n", i, vote)
	}
	expectedAverage := sumVotes / float64(n)
	fmt.Printf("\nExpected network average: %.4f (majority: %d)\n\n", expectedAverage, func() int {
		if expectedAverage > 0.5 {
			return 1
		}
		return 0
	}())

	// Build topology
	for i := range n {
		next := (i + 1) % n

		// Channel i -> next
		chForward := make(chan float64, 1)
		nodes[i].out = append(nodes[i].out, chForward)
		nodes[next].in = append(nodes[next].in, chForward)

		// Channel next -> i
		chBackward := make(chan float64, 1)
		nodes[next].out = append(nodes[next].out, chBackward)
		nodes[i].in = append(nodes[i].in, chBackward)
	}

	// Launch goroutines
	wg.Add(n)
	for _, nd := range nodes {
		go nd.start()
	}

	wg.Wait()
	fmt.Println("\nDone")
}
