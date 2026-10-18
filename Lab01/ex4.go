package main
import (
	"fmt"
	"os"
	"bufio"
	"strconv"
)

func main() {
	slice := []float32{23.2, 123.123, 1.1, 0.00001, 125567.5}
	fmt.Println("slice initial: ", slice)

	fmt.Print("Introdu index: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()
	index, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("Err la conversie")
	} else if index > len(slice - 1) {
		fmt.Println("Invalid index")
	} else {
		slice = append(slice[:index], slice[index+1:]...)
		fmt.Println("New slice: ", slice)
	}
}
