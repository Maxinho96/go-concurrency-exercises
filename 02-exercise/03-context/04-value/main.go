package main

import (
	"context"
	"fmt"
)

type user string

func main() {
	ctx, done := context.WithCancel(context.Background())
	defer done()
	processRequest(ctx, "jane")
}

func processRequest(ctx context.Context, userid string) {
	// TODO: send userID information to checkMemberShip through context for
	// map lookup.
	ctx = context.WithValue(ctx, user("jane"), true)
	ch := checkMemberShip(ctx)
	status := <-ch
	fmt.Printf("membership status of userid : %s : %v\n", userid, status)
}

// checkMemberShip - takes context as input.
// extracts the user id information from context.
// spins a goroutine to do map lookup
// sends the result on the returned channel.
func checkMemberShip(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	go func() {
		defer close(ch)
		// do some database lookup
		status := ctx.Value(user("jane")).(bool)
		ch <- status
	}()
	return ch
}
