package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs

		switch sig {
		case syscall.SIGINT:
			fmt.Println("\nreceived SIGINT")
		case syscall.SIGTERM:
			fmt.Println("\nreceived SIGTERM")
		}

		done <- true
	}()

	fmt.Println("awaiting sinal")
	<-done
	fmt.Println("exiting")
	os.Exit(1)
}
