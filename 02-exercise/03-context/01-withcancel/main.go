package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

func main() {

	// TODO: generator -  generates integers in a separate goroutine and
	// sends them to the returned channel.
	// The callers of gen need to cancel the goroutine once
	// they consume 5th integer value
	// so that internal goroutine
	// started by gen is not leaked.
	generator := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for i := 0; true; i++ {
				select {
				case <-ctx.Done():
					fmt.Println("closed")
					return
				case ch <- i:
				}
			}
		}()
		return ch
	}

	// Create a context that is cancellable.
	ctx, cancel := context.WithCancel(context.Background())
	ch := generator(ctx)
	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)
	}
	cancel()
	time.Sleep(time.Millisecond * 10)
}
