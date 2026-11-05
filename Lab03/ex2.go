package main
import (
	"fmt"
	"sort"
	"math/rand"
)

type LamportClock struct { counter int }

func (l *LamportClock ) Tick () {
	l.counter++
}

func (l *LamportClock ) Send () int {
	l.counter++
	return l.counter
}

func (l *LamportClock ) Receive (ts int) {
	if ts > l.counter {
		l.counter = ts
	}
	l.counter++
}

type Message struct { sender , ts int }

type Event struct {
	process int
	internal_counter int
	kind string
}

func process(id int, send chan<- Message, recv <-chan Message, main chan<- Event){
	clock := LamportClock {0}

	for range(1 + rand.Intn(4)) {
		clock.Tick()
		main <- Event {id, clock.counter, "internal"}
	}

	send <- Message {id, clock.Send()}
	main <- Event {id, clock.counter, "send"}

	for range(1 + rand.Intn(4)) {
		clock.Tick()
		main <- Event {id, clock.counter, "internal"}
	}

	msg := <- recv
	clock.Receive(msg.ts)
	main <- Event {id, clock.counter, "recv"}
	
	close(main)
}

func main() {
	ch12 := make(chan Message, 1)
	ch23 := make(chan Message, 1)
	ch31 := make(chan Message, 1)
	chMain1 := make(chan Event, 10)
	chMain2 := make(chan Event, 10)
	chMain3 := make(chan Event, 10)

	go process(1, ch12, ch31, chMain1)
	go process(2, ch23, ch12, chMain2)
	go process(3, ch31, ch23, chMain3)

	events := make([]Event, 0, 30)

	for event := range chMain1 { events = append(events, event) }
	for event := range chMain2 { events = append(events, event) }
	for event := range chMain3 { events = append(events, event) }
	
	sort.Slice(events, func(i, j int) bool {
		return less(events[i], events[j])
	})

	for i, event := range events {
		fmt.Printf("%v ", event)
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}
}

func less(a, b Event) bool {
	if (a.internal_counter == b.internal_counter) {
		return a.process < b.process
	}
	return a.internal_counter < b.internal_counter
}
