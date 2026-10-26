package main
import "fmt"

func main() {
	var ch = make(chan string)

	go sender(ch)
	fmt.Println("Success on starting sender from main")

	for msg := range ch {
		fmt.Println("Success on recv message: ", msg)
	}
	// msg, ok := <- ch
	// for ok != false {
	// 	fmt.Println("Success on recv message: ", msg)
	//
	// 	msg, ok = <- ch
	// }

	fmt.Println("End inside main")
}

func sender(ch chan<- string) {
	fmt.Println("Start inside sender")

	fmt.Println("Send 1")
	ch <- "Mesaj 1"
	fmt.Println("1 sent\nSend 2")
	ch <- "Mesaj 2"
	fmt.Println("2 sent\nClose ch")
	close(ch)

	fmt.Println("End inside sender")
}
