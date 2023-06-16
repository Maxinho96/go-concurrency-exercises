package main

import "fmt"

// TODO: Implement relaying of message with Channel Direction

func genMsg(out chan<- string) {
	// send message on ch1
	out <- "Message"
}

func relayMsg(in <-chan string, out chan<- string) {
	// recv message on ch1
	message := <-in
	// send it on ch2
	out <- message
}

func main() {
	// create ch1 and ch2
	ch1 := make(chan string)
	ch2 := make(chan string)
	// spine goroutine genMsg and relayMsg
	go genMsg(ch1)
	go relayMsg(ch1, ch2)
	// recv message on ch2
	fmt.Println(<-ch2)
}
