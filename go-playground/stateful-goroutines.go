package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type readOp struct {
	key int
	res chan int
}

type writeOp struct {
	key int
	val int
	res chan bool
}

func main() {
	var readOps uint64
	var writeOps uint64

	reads := make(chan readOp)
	writes := make(chan writeOp)

	go func() {
		var state = make(map[int]int)
		for {
			select {
			case read := <-reads:
				read.res <- state[read.key]
			case write := <-writes:
				state[write.key] = write.val
				write.res <- true
			}
		}
	}()

	for w := 0; w < 100; w++ {
		go func(i int) {
			write := writeOp{
				key: i,
				val: i * 2,
				res: make(chan bool),
			}
			writes <- write
			<-write.res
			atomic.AddUint64(&writeOps, 1)
			time.Sleep(time.Millisecond)
		}(w)
	}

	for r := 0; r < 100; r++ {
		go func(i int) {
			read := readOp{
				key: i,
				res: make(chan int),
			}
			reads <- read
			res := <-read.res
			fmt.Println("reading -> key: ", read.key, "val: ", res)
			atomic.AddUint64(&readOps, 1)
			time.Sleep(time.Millisecond)
		}(r)
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)
}
