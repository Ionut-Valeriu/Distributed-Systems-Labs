package main
import "fmt"

func main()  {
	sir := [...]int{1,4,7,12,11,145,25,16}
	fmt.Println("Sir initial: ", sir)

	var slice = make([]int, 0, len(sir))

	for _, nr := range sir  {
		if nr % 2 == 0 {
			slice = append(slice, nr)
		}
	}

	fmt.Println("Slice: ", slice)
}
