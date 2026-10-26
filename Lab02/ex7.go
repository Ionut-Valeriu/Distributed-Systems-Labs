package main
import (
	"fmt"
	"sync"
	"math/rand"
	"strconv"
)

type coord struct {
	x int
	y int
}

func (this coord) make_diff (other coord) coord {
	x := this.x - other.x
	y := this.y - other.y
	return coord{x, y}
}

func (this coord) show() string {
	return "{x: " + strconv.Itoa(this.x) + " ; y: " + strconv.Itoa(this.y) + "}"
}

func main() {
	raport := make(chan string)
	ch1to2 := make(chan coord, 1)
	ch2to1 := make(chan coord, 1)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go robo(ch1to2, ch2to1, raport, &wg)
	go robo(ch2to1, ch1to2, raport, &wg)

	go func() {
		wg.Wait()
		close(raport)
	}()

	for msg := range raport {
		fmt.Println(msg)
	}
}

func robo(in <-chan coord, out chan<- coord, raport chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	x := rand.Intn(50)
	y := rand.Intn(50)
	my_coords := coord{x, y}
	intial := my_coords

	out <- my_coords
	for other_coords := range in {
		diff := my_coords.make_diff(other_coords)
		if diff.x == 0 && diff.y == 0 {
			close(out)
			break
		}
		move(&my_coords, diff)
		out <- my_coords
	}

	raport <- "I started at: " + intial.show() + " and end at: " + my_coords.show();
}

func move(my_coords *coord, diff coord) {
	if (diff.x != 0 && diff.x != -1) {
		if (diff.x < 0) {
			(*my_coords).x ++
		} else {
			(*my_coords).x --
		}
	} else if diff.y != -1 {
		if (diff.y < 0) {
			(*my_coords).y ++
		} else {
			(*my_coords).y --
		}
	}
}
