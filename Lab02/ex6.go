package main
import (
	"fmt"
	"time"
)

// observer
func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go robo(ch1)
	go robo(ch2)

	select {
	case msg := <-ch1:
		fmt.Println("Robo 1: ", msg)
	case msg := <-ch2:
		fmt.Println("Robo 2: ", msg)
	case <-time.After(1000 * time.Millisecond):
		fmt.Println("No response")
	}
}

func robo(ch chan<- string) {
	ch <- "I finished!"
	close(ch)
}
