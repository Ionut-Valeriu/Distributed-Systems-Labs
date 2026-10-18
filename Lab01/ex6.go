package main
import "fmt"

func main() {
	myMap := map[string]int{
		"Andrei ": 190,
		"Maria": 165,
		"Ion": 182,
		"Elena": 172,
		"George ": 178,
		"Ana": 165,
	}
	
	names := make([]string, 0, len(myMap))

	for name, _ := range myMap {
		names = append(names, name)
	}

	for i := range names {
		for j := i + 1; j < len(names); j++ {
			if myMap[names[i]] < myMap[names[j]] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	fmt.Println(names)
}
