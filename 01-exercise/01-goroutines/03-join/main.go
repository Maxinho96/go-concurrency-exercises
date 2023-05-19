package main

import (
	"fmt"
	"sync"
)

func main() {
	//TODO: modify the program
	// to print the value as 1
	// deterministically.

	var wg sync.WaitGroup
	wg.Add(1)

	var data int

	go func() {
		defer wg.Done()
		data++
	}()

	wg.Wait()

	fmt.Printf("the value of data is %v\n", data)

	fmt.Println("Done..")
}
