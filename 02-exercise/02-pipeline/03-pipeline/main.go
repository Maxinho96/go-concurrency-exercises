// generator() -> square() ->
//
//															-> merge -> print
//	            -> square() ->
package main

import (
	"fmt"
	"sync"
)

func generator(done <-chan any, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}

		}
	}()
	return out
}

func square(done <-chan any, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan any, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	done := make(chan any)
	in := generator(done, 2, 3)

	c1 := square(done, in)
	c2 := square(done, in)

	out := merge(done, c1, c2)

	// TODO: cancel goroutines after receiving one value.

	fmt.Println(<-out)
	close(done)
}
