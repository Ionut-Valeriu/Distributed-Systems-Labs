package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Message reprezinta un mesaj transmis intre noduri in Runda 1 sau 2
type Message struct {
	SenderID int
	Round    int // 1 pentru starea curenta, 2 pentru valoarea Tie
	Value    int // 0 sau 1
}

type Node struct {
	id          int
	x           int // starea curenta x_i(t)
	s           int
	n           int  // numarul de noduri din sistem
	isByzantine bool // comportamentul la iteratia curenta
	inbox       chan Message
}

func (node *Node) BcastState(nodes []Node) {
	for i := range nodes {
		nodes[i].inbox <- Message{
			SenderID: node.id,
			Round:    1,
			Value:    node.x,
		}
	}
}

func (node *Node) MaliciousBcastState(nodes []Node) {
	for i := range nodes {
		nodes[i].inbox <- Message{
			SenderID: node.id,
			Round:    1,
			Value:    rand.Intn(2),
		}
	}
}

func (node *Node) RecvState() (int, int) {
	ones := 0
	for range node.n {
		msg := <-node.inbox
		if msg.Round == 1 && msg.Value == 1 {
			ones++
		}
	}
	if ones >= node.n/2+1 {
		return 1, ones
	}
	return 0, node.n - ones
}

func (node *Node) BcastTie(nodes []Node, majority int) {
	for i := range nodes {
		if i == node.id {
			continue
		}
		nodes[i].inbox <- Message{
			SenderID: node.id,
			Round:    2,
			Value:    majority,
		}
	}
}

func (node *Node) MaliciousBcastTie(nodes []Node) {
	for i := range nodes {
		if i == node.id {
			continue
		}
		nodes[i].inbox <- Message{
			SenderID: node.id,
			Round:    2,
			Value:    rand.Intn(2),
		}
	}
}

func (node *Node) RecvTie(leaderID int, majority int) int {
	if node.id == leaderID {
		return majority
	}
	msg := <-node.inbox
	return msg.Value
}

func main() {
	nodes, s := initNodes()
	if nodes == nil {
		panic(1)
	}

	n := len(nodes)
	majorities := make([]int, n)
	mults := make([]int, n)

	byzantineNodes := make([]int, 0, s)
	initX := make([]int, n)
	for i := range nodes {
		initX[i] = nodes[i].x
		if nodes[i].isByzantine {
			byzantineNodes = append(byzantineNodes, nodes[i].id)
		}
	}
	fmt.Printf("Noduri bizantine: %v |   x(0): %v\n", byzantineNodes, initX)

	t := 0
	for {
		var wg sync.WaitGroup

		// Runda 1 (Agregarea voturilor): broadcast starea curenta x_i(t)
		for i := range nodes {
			wg.Add(1)
			go func(node *Node) {
				defer wg.Done()
				if node.isByzantine {
					node.MaliciousBcastState(nodes)
				} else {
					node.BcastState(nodes)
				}
			}(&nodes[i])
		}
		wg.Wait()

		// Runda 1: primire mesaje si calcul Majority(t) si mult(t)
		for i := range nodes {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				majorities[idx], mults[idx] = nodes[idx].RecvState()
			}(i)
		}
		wg.Wait()

		// Runda 2 (Spargerea simetriei prin Leader/King): broadcast Tie de la liderul P_t
		leaderID := t % n
		if nodes[leaderID].isByzantine {
			nodes[leaderID].MaliciousBcastTie(nodes)
		} else {
			nodes[leaderID].BcastTie(nodes, majorities[leaderID])
		}

		// Runda 2: receptie Tie si actualizarea starii x_i(t+1)
		for i := range nodes {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				tie := nodes[idx].RecvTie(leaderID, majorities[idx])
				if !nodes[idx].isByzantine {
					if mults[idx] > nodes[idx].n/2+nodes[idx].s {
						nodes[idx].x = majorities[idx]
					} else {
						nodes[idx].x = tie
					}
				}
			}(i)
		}
		wg.Wait()

		xs := make([]int, n)
		for i := range nodes {
			xs[i] = nodes[i].x
		}
		leaderRole := "corect  "
		if nodes[leaderID].isByzantine {
			leaderRole = "bizantin"
		}
		fmt.Printf("Iter %d | Lider: %d (%s) | x: %v\n", t, leaderID, leaderRole, xs)

		// Conditia de oprire
		if t > s+1 {
			break
		}

		t++
	}

	for i := range nodes {
		if !nodes[i].isByzantine {
			fmt.Printf("Consens: %d\n", nodes[i].x)
			break
		}
	}
}

func initNodes() ([]Node, int) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil, 0
	}

	firstLine := strings.TrimSpace(scanner.Text())
	var n, s int
	_, err := fmt.Sscan(firstLine, &n, &s)
	if err != nil {
		// Daca prima linie este un nume de fisier (stil Lab06), citim din fisier
		file, ferr := os.Open(firstLine)
		if ferr != nil {
			fmt.Println("File not found")
			return nil, 0
		}
		defer file.Close()

		scanner = bufio.NewScanner(file)
		if !scanner.Scan() {
			return nil, 0
		}
		_, err = fmt.Sscan(scanner.Text(), &n, &s)
		if err != nil {
			return nil, 0
		}
	}

	if n <= 0 || s < 0 || n <= 4*s {
		fmt.Println("Invalid parameters: must satisfy n > 4*s and s >= 0")
		return nil, 0
	}

	values, err := readValues(scanner, n)
	if err != nil {
		return nil, 0
	}

	if err := scanner.Err(); err != nil {
		return nil, 0
	}

	nodes := make([]Node, n)
	for i := range n {
		nodes[i] = Node{
			id:    i,
			x:     values[i],
			s:     s,
			n:     n,
			inbox: make(chan Message, n),
		}
	}

	// Alegem aleator s noduri bizantine distincte
	perm := rand.Perm(n)
	for i := range s {
		nodes[perm[i]].isByzantine = true
	}

	return nodes, s
}

func readValues(scanner *bufio.Scanner, n int) ([]int, error) {
	if !scanner.Scan() {
		return nil, fmt.Errorf("missing initial values line")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d values, got %d", n, len(fields))
	}
	values := make([]int, n)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil || (val != 0 && val != 1) {
			return nil, fmt.Errorf("invalid binary value %q", f)
		}
		values[i] = val
	}
	return values, nil
}
