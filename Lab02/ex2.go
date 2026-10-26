package main
import (
	"fmt"
	"sync"
)

func main() {
	var ch = make(chan string)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go receiver(ch, &wg)
	fmt.Println("Success on starting receiver from main")

	ch <- "Salut"
	fmt.Println("Success on send message")
	close(ch)

	wg.Wait()
	fmt.Println("End inside main")
}

func receiver(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Start inside receiver")
	fmt.Println(<- ch)
	fmt.Println("End inside receiver")
}
