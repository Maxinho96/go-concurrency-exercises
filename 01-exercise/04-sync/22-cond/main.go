package main

import (
	"fmt"
	"sync"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	c := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.L.Unlock()
		c.L.Lock()

		//TODO: suspend goroutine until sharedRsc is populated.

		for len(sharedRsc) < 1 {
			c.Wait()
		}

		fmt.Println(sharedRsc["rsc1"])
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.L.Unlock()
		c.L.Lock()

		//TODO: suspend goroutine until sharedRsc is populated.

		for len(sharedRsc) < 2 {
			c.Wait()
		}

		fmt.Println(sharedRsc["rsc2"])
	}()

	c.L.Lock()
	// writes changes to sharedRsc
	sharedRsc["rsc1"] = "foo"
	sharedRsc["rsc2"] = "bar"
	c.Broadcast()
	c.L.Unlock()

	wg.Wait()
}
