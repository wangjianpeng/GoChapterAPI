package mastergo

import (
	"fmt"
	"time"
)

func DoRunChannel() {
	ch1 := make(chan int)
	go doSum(ch1, "first")
	go doSum(ch1, "second")
	// i, j := <-ch1, <-ch1
	i := <-ch1
	j := <-ch1
	fmt.Println("end main.", i, j)
}

func doSum(ch chan<- int, name string) {
	time.Sleep(2 * time.Second)
	ch <- 10
	fmt.Println("sun name: ", name)
}
