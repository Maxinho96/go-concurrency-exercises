package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cond.L.Unlock()

		//TODO: suspend goroutine until sharedRsc is populated.
		cond.L.Lock()
		for len(sharedRsc) == 0 {
			cond.Wait()
		}

		fmt.Println(sharedRsc["rsc1"])
	}()

	time.Sleep(time.Second)
	// writes changes to sharedRsc
	cond.L.Lock()
	sharedRsc["rsc1"] = "foo"
	cond.Signal()
	cond.L.Unlock()

	wg.Wait()
}
