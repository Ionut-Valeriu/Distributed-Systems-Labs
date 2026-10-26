package main
import "fmt"

func main() {
	ch := make(chan string)
	done := make (chan bool)

	fmt.Println("Recv start")
	go recv(ch, done)

	ch <- "msg 1"
	ch <- "msg 12"
	ch <- "msg 13"
	ch <- "msg 14"
	close(ch)

	for range done {}
	fmt.Println("Recv is done. (main)")
}

func recv(ch <-chan string, done chan<- bool) {
	for msg := range ch {
		fmt.Println("Recv: ", msg)
	}

	fmt.Println("Recv is done. (recv)")
	done <- true
	close (done)
}
