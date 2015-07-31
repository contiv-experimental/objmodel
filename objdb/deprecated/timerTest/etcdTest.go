// This package is to test some strange timer behavior seen in go runtime
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/contiv/symphony/pkg/confStore"
	api "github.com/contiv/symphony/pkg/confStore/confStoreApi"
)

func main() {
	cStore := confStore.NewConfStore()

	hostName, _ := os.Hostname()
	myId := hostName
	fmt.Printf("My Identifier: %s\n", myId)

	// Create a lock
	lock, err := cStore.NewLock("master", myId, 10)

	// Acquire the lock
	err = lock.Acquire(0)
	if err != nil {
		fmt.Printf("Error acquiring lock\n")
	}

	cnt := 0
	for {
		select {
		case event := <-lock.EventChan():
			fmt.Printf("Event on Lock1: %+v\n\n", event)
			if event.EventType == api.LockAcquired {
				fmt.Printf("Master lock acquired by %s\n", myId)
			}
		case <-time.After(time.Second * time.Duration(10)):
			// At this point, lock should be holding the lock
			if lock.IsAcquired() {
				fmt.Printf("Lock is still acquired. I'm the master %s\n", myId)
				// Release the lock once in a while
				/*if ((cnt % 5) == 4) {
				    fmt.Printf("Releasing the master lock: %s\n", myId)

				    // Release the lock
				    lock.Release()

				    // Wait a little
				    time.Sleep(time.Second * 3)

				    // Create a brand new lock
				    lock, _ = cStore.NewLock("master", myId, 10)

				    // reAcquire the lock
				    err = lock.Acquire(0)
				    if (err != nil) {
				        fmt.Printf("Error acquiring lock")
				    }
				}*/
			} else {
				fmt.Printf("Master lock is held by %s\n", lock.GetHolder())
			}

			cnt++
		}
	}
}
