package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Pair struct {
	Id  int
	Msg int
}

// Memory este o multime de token -uri (reuniunea se face prin inserare).
type Memory map[Pair]struct{}
type Node struct {
	ID      int
	M       Memory
	In      []chan Memory
	Out     []chan Memory
	Crashed bool
}

func (n *Node) Crash() {
	n.Crashed = true
	for _, ch := range n.Out {
		close(ch)
	}
}

func (n *Node) Bcast() {
	snapshot := make(Memory, len(n.M))
	for k := range n.M {
		snapshot[k] = struct{}{}
	}
	for _, ch := range n.Out {
		select {
		case ch <- snapshot: // if a neighbor crash, this is not blocking
		default:
		}
	}
}

func (n *Node) Recv() {
	for _, ch := range n.In {
		if msg, ok := <-ch; ok {
			for k := range msg {
				n.M[k] = struct{}{}
			}
		}
	}
}

func main() {
	var filename string
	_, err := fmt.Scanln(&filename)
	if err != nil {
		panic(1)
	}

	nodes, s, diam := initNodes(filename)
	if nodes == nil {
		panic(1)
	}

	crashedCount := 0
	t := 0
	for {
		// Decide random if a node crashes in this iteration (none in iteration 0, at most s crashes)
		if t > 0 && crashedCount < s && rand.Intn(2) == 1 {
			alive := make([]int, 0, len(nodes))
			for i := range nodes {
				if !nodes[i].Crashed {
					alive = append(alive, i)
				}
			}
			if len(alive) > 1 {
				victim := alive[rand.Intn(len(alive))]
				nodes[victim].Crash()
				fmt.Printf("Iteratia %d: nodul %d a intrat in crash\n", t, nodes[victim].ID)
				crashedCount++
			}
		}

		// 1. Broadcast M_i(t) to all neighbors
		var wg sync.WaitGroup
		for i := range nodes {
			if nodes[i].Crashed {
				continue
			}
			wg.Add(1)
			go func(node *Node) {
				defer wg.Done()
				node.Bcast()
			}(&nodes[i])
		}
		wg.Wait()

		// 2-3. Receive from all neighbors and merge into M_i(t+1)
		for i := range nodes {
			if nodes[i].Crashed {
				continue
			}
			wg.Add(1)
			go func(node *Node) {
				defer wg.Done()
				node.Recv()
			}(&nodes[i])
		}
		wg.Wait()

		// 4. Stopping condition
		if t > (s+1)*diam {
			break
		}

		// 5. Next iteration
		t++
	}

	fmt.Printf("Diametru: %d\n", diam)
	fmt.Printf("Iteratie oprire: %d\n", t)

	for i := range nodes {
		if !nodes[i].Crashed {
			pairs := make([]Pair, 0, len(nodes[i].M))
			for p := range nodes[i].M {
				pairs = append(pairs, p)
			}
			sort.Slice(pairs, func(a, b int) bool {
				return pairs[a].Id < pairs[b].Id
			})
			for j, p := range pairs {
				if j > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("(%d, %d)", p.Id, p.Msg)
			}
			fmt.Println()
			break
		}
	}
}

func initNodes(filename string) ([]Node, int, int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("File not found")
		return nil, 0, -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 1st line: n and s
	if !scanner.Scan() {
		return nil, 0, -1
	}
	var n, s int
	_, err = fmt.Sscan(scanner.Text(), &n, &s)
	if err != nil {
		return nil, 0, -1
	}

	// 2nd line: n token values
	tokens, err := readTokens(scanner, n)
	if err != nil {
		return nil, 0, -1
	}

	// next n lines: adjacency lists
	adj, err := readAdj(scanner, n)
	if err != nil {
		return nil, 0, -1
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, -1
	}

	diam := diametru(adj)
	if diam == -1 {
		fmt.Println("Graph is not strongly connected")
		return nil, 0, -1
	}

	nodes := make([]Node, n)
	for i := range n {
		nodes[i] = Node{
			ID: i,
			M: Memory{
				Pair{Id: i, Msg: tokens[i]}: {},
			},
		}
	}

	// connect directed edges i -> neighbor via channels
	for i := range n {
		for _, neighbor := range adj[i] {
			ch := make(chan Memory, 1)
			nodes[i].Out = append(nodes[i].Out, ch)
			nodes[neighbor].In = append(nodes[neighbor].In, ch)
		}
	}

	return nodes, s, diam
}

func readTokens(scanner *bufio.Scanner, n int) ([]int, error) {
	if !scanner.Scan() {
		return nil, fmt.Errorf("missing tokens line")
	}
	tokenFields := strings.Fields(scanner.Text())
	if len(tokenFields) != n {
		return nil, fmt.Errorf("expected %d tokens, got %d", n, len(tokenFields))
	}
	tokens := make([]int, n)
	for i, f := range tokenFields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		tokens[i] = val
	}
	return tokens, nil
}

func readAdj(scanner *bufio.Scanner, n int) ([][]int, error) {
	adj := make([][]int, n)
	for i := range n {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing adjacency list for node %d", i)
		}
		fields := strings.Fields(scanner.Text())
		adj[i] = make([]int, 0, len(fields))
		for _, f := range fields {
			neighbor, err := strconv.Atoi(f)
			if err != nil || neighbor < 0 || neighbor >= n {
				return nil, fmt.Errorf("invalid neighbor %q for node %d", f, i)
			}
			adj[i] = append(adj[i], neighbor)
		}
	}
	return adj, nil
}

func diametru(adj [][]int) int {
	n := len(adj)
	diam := 0
	for src := 0; src < n; src++ {
		// BFS din nodul src
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		dist[src] = 0
		queue := []int{src}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		for _, d := range dist {
			if d == -1 {
				return -1 // avem un nod la care nu putem ajunge din src
			}
			if d > diam {
				diam = d
			}
		}
	}
	return diam
}
