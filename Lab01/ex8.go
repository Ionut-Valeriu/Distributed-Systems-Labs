package main
import "fmt"

func main() {
	jobers := map[string]int{
		"Ana Popescu": 1,
		"Mihai Ionescu": 2,
		"Bogdan Radu": 1,
		"Elena Dumitrescu": 3,
		"Cristian Vasile": 2,
		"Ioana Enache": 1,
	}

	departs := make(map[int][]string)

	for name, id := range jobers {
		departs[id] = append(departs[id], name)
	}

	for id, list := range departs {
		fmt.Println("Depart", id)
		for index := range list {
			fmt.Println("\t", list[index])
		}
	}
}
