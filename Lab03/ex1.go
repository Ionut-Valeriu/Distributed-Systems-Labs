package main
import (
	"fmt"
	"time"
	"math/rand"
)

type Clock struct { offset time.Duration }

func (c *Clock) Now () time.Time { return time.Now().Add(c.offset) }

func estimate (t1, t2, t3, t4 time.Time) (theta, delta time.Duration) {
	theta = (t2.Sub(t1) + t3.Sub(t4)) / 2
	delta = t4.Sub(t1) - t3.Sub(t2)
	return theta, delta
}

func measure(c_clock, s_clock *Clock) (theta, delta time.Duration) {
	t1 := c_clock.Now()
	time.Sleep( 3 * time.Millisecond + time.Duration(rand.Intn(10)) * time.Millisecond )

	t2 := s_clock.Now()
	time.Sleep( 5 * time.Millisecond + time.Duration(rand.Intn(20)) * time.Millisecond )

	t3 := s_clock.Now()
	time.Sleep( 3 * time.Millisecond + time.Duration(rand.Intn(10)) * time.Millisecond )

	t4 := c_clock.Now()
	return estimate(t1, t2, t3, t4)
}

func main() {
	fmt.Println("A.");
	client := &Clock{200 * time.Millisecond}
	server := &Clock{0 * time.Millisecond}
	fmt.Println("Initial offset: ", client.offset)

	theta, delta := measure(client, server)
	client.offset += theta

	fmt.Println("Thetha:", theta, "; Delta/2:", delta/2, "\nNew client:", client.offset);

	fmt.Println("B.");
	var minT, minD time.Duration
	client = &Clock{200 * time.Millisecond}
	for i := range(10) {
		t, d := measure(client, server)
		if i == 0 || minD > d {
			minT, minD = t, d
		}
	}
	fmt.Println("Minim:  ", minT, minD)
	fmt.Println("Point A:", theta, delta)
}
