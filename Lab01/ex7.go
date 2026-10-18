package main
import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	myMap := make(map[string]int)
	file, err := os.Open("file7.txt")

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) 

	for scanner.Scan() {
		myMap[scanner.Text()] ++;
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println("Eroare la citire: ", err)
	}

	fmt.Println(myMap)
}
