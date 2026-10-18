package main
import "fmt"

func main() {
	myArr := [4]int{3,2,4,1}

	fmt.Printf("initial arr: %v\n", myArr)

	for i := range(len(myArr)) {
		for j := i+1; j < len(myArr); j++ {
			if myArr[i] > myArr[j]{
				aux := myArr[i]
				myArr[i] = myArr[j]
				myArr[j] = aux
			}
		}
	}

	fmt.Printf("sorted: %v", myArr)
}
