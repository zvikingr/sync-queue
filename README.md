### SyncQueue
A simple, goroutine-safe, sync queue for Go

### Installation

    go get github.com/okayping/sync-queue

### Example:
```golang
package main

import (
	"log"
	"time"

	queue "github.com/okayping/sync-queue"
)

func main() {
	q := queue.NewSyncQueue()

	go func() {
		var i = 0
		for {
			time.Sleep(time.Millisecond * 500)
			i += 1
			q.Push(i)
		}
	}()

	time.Sleep(time.Second)

	for {
		v := q.Pop()
		if v == nil {
			break
		}

		log.Println("get value:", v.(int))
	}
}
```
