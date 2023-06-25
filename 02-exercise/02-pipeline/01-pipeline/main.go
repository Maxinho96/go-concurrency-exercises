package main

import "fmt"

// TODO: Build a Pipeline
// generator() -> square() -> print

// generator - convertes a list of integers to a channel
func generator(nums ...int) <-chan int {
	out_ch := make(chan int)
	go func() {
		defer close(out_ch)
		for _, num := range nums {
			fmt.Println("Generating", num)
			out_ch <- num
		}
	}()
	return out_ch
}

// square - receive on inbound channel
// square the number
// output on outbound channel
func square(in_ch <-chan int) <-chan int {
	out_ch := make(chan int)
	go func() {
		defer close(out_ch)
		for num := range in_ch {
			fmt.Println("Squaring", num)
			out_ch <- num * num
		}
	}()
	return out_ch
}

func main() {
	// set up the pipeline
	out_ch := square(generator(1, 2, 3))
	// run the last stage of pipeline
	// receive the values from square stage
	// print each one, until channel is closed.
	for num := range out_ch {
		fmt.Println("Result:", num)
	}
}
