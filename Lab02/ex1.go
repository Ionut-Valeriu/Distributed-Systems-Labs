package main
import "fmt"

func main() {
	var ch = make(chan string)

	go sender(ch)
	fmt.Println("Success on starting sender from main")

	fmt.Println(<- ch)
	fmt.Println("Success on recv message")

	fmt.Println("End inside main")
}

func sender(ch chan string) {
	fmt.Println("Start inside sender")

	ch <- "Salut"
	close(ch)

	fmt.Println("End inside sender")
}
