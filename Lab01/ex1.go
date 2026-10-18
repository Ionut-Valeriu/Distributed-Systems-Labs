package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func main() {
	var number int

	fmt.Printf("Scire nr: ");
	_, err := fmt.Scanf("%d", &number)

	if (err != nil) {
		fmt.Printf("nr necitit")
	} else {
		if (isPrime(number)) {
			fmt.Printf("%d is prime\n", number)
		} else {
			fmt.Printf("%d is not prime\n", number)
		}
	}
}

func isPrime(x int) bool {
	if (x < 2) {
		return false
	}
	if (x == 2) {
		return true
	}
	if (x % 2 == 0) {
		return false
	}
	for i := 3; i * i <= x ; i++ {
		if x % i == 0 {
			return false
		}
	}
	return true
}

func main2() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	v := scanner.Text()
	x, err := strconv.Atoi(v)

	if err != nil {
		fmt.Println("Err la conversie")
	}

	fmt.Printf("%T", x)
}
