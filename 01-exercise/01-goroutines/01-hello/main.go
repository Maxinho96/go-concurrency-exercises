package main

import (
	"fmt"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	// Direct call
	fun("direct call")

	// TODO: write goroutine with different variants for function call.

	// goroutine function call
	go fun("goroutine function call")

	// goroutine with anonymous function
	go func() {
		fun("anonymous function call")
	}()

	// goroutine with function value call
	function_value := fun
	go function_value("function value call")

	// wait for goroutines to end
	time.Sleep(1 * time.Second)

	fmt.Println("done..")
}
