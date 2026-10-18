package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	input := scanner.Text()
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Err la conversie")
	}


	scanner.Scan()
	input = scanner.Text()
	s, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Err la conversie")
	}

	myMap := make(map[int]int)
	for i := range n {
		scanner.Scan()
		input = scanner.Text()
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Err la conversie")
		}

		myMap[number] = i;
	
		ind, ok := myMap[s - number]
		if ok == true {
			fmt.Println(ind, i)
		}
	}
}
