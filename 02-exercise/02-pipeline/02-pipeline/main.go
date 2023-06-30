// generator() -> square() -> print

package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	// Implement fan-in
	// merge a list of channels to a single channel
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range cs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for n := range ch {
				out <- n
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := generator(2, 3)

	// TODO: fan out square stage to run two instances.
	out1 := square(in)
	out2 := square(in)

	// TODO: fan in the results of square stages.
	out := merge(out1, out2)
	for i := range out {
		fmt.Println(i)
	}
}
