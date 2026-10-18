package main
import "fmt"

func main() {
	slice := []string{"Laboratorul", "de", "sisteme", "distribuite"}

	for i := range len(slice) / 2 {
		slice[i], slice[len(slice) - 1 - i] = slice[len(slice) - 1 - i], slice[i]
	}

	fmt.Println(slice)
}
