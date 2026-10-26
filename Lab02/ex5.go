package main
import "fmt"

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go send(ch1)
	go recv(ch1, ch2)

	for msg := range ch2 {
		fmt.Println(msg)
	}
}

func send(ch chan<- string) {
	ch <- "Salut"
	close(ch)
}

func recv(ch1 <-chan string, ch2 chan<- string) {
	for msg := range ch1 {
		ch2 <- msg
	}
	close(ch2)
}
