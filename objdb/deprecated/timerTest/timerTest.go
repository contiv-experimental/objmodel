// This package is to test some strange timer behavior seen in go runtime
package main

import (
	"fmt"
	"time"

	"github.com/contiv/symphony/pkg/confStore"
	api "github.com/contiv/symphony/pkg/confStore/confStoreApi"
)

// New confstore
var cStore api.ConfStorePlugin

func init() {
	cStore = confStore.NewConfStore()

	// First test
	TestLockAcquireRelease()

	// Second test
	TestLockAcquireTimeout()
}

func TestLockAcquireRelease() {
	// Create a lock
	lock1, err := cStore.NewLock("master", "hostname1", 10)
	lock2, err := cStore.NewLock("master", "hostname2", 10)

	// Acquire the lock
	err = lock1.Acquire(0)
	if err != nil {
		fmt.Printf("Error acquiring lock1")
	}
	err = lock2.Acquire(0)
	if err != nil {
		fmt.Printf("Error acquiring lock2")
	}

	cnt := 1
	for {
		select {
		case event := <-lock1.EventChan():
			fmt.Printf("Event on Lock1: %+v\n\n", event)
			if event.EventType == api.LockAcquired {
				fmt.Printf("Master lock acquired by Lock1\n")
			}
		case event := <-lock2.EventChan():
			fmt.Printf("Event on Lock2: %+v\n\n", event)
			if event.EventType == api.LockAcquired {
				fmt.Printf("Master lock acquired by Lock2\n")
			}
		case <-time.After(time.Second * time.Duration(30)):
			if cnt == 1 {
				fmt.Printf("10sec timer. releasing Lock1\n\n")
				// At this point, lock1 should be holding the lock
				if !lock1.IsAcquired() {
					fmt.Printf("Lock1 failed to acquire lock\n\n")
				}
				lock1.Release()
				cnt++
			} else {
				fmt.Printf("20sec timer. releasing Lock2\n\n")

				// At this point, lock1 should be holding the lock
				if !lock2.IsAcquired() {
					fmt.Printf("Lock2 failed to acquire lock\n\n")
				}

				// we are done with the test
				lock2.Release()

				return
			}
		}
	}
}

func TestLockAcquireTimeout() {
	fmt.Printf("\n\n\n\n\n\n =========================================================== \n\n\n\n\n")
	// Create a lock
	lock1, err := cStore.NewLock("master", "hostname1", 10)
	lock2, err := cStore.NewLock("master", "hostname2", 10)

	// Acquire the lock
	err = lock1.Acquire(0)
	if err != nil {
		fmt.Printf("Error acquiring lock1")
	}
	err = lock2.Acquire(20)
	if err != nil {
		fmt.Printf("Error acquiring lock2")
	}

	for {
		select {
		case event := <-lock1.EventChan():
			fmt.Printf("Event on Lock1: %+v\n\n", event)
			if event.EventType == api.LockAcquired {
				fmt.Printf("Master lock acquired by Lock1\n")
			}
		case event := <-lock2.EventChan():
			fmt.Printf("Event on Lock2: %+v\n\n", event)
			if event.EventType != api.LockAcquireTimeout {
				fmt.Printf("Invalid event on Lock2\n")
			} else {
				fmt.Printf("Lock2 timeout as expected\n")
			}
		case <-time.After(time.Second * time.Duration(40)):
			fmt.Printf("40sec timer. releasing Lock1\n\n")
			// At this point, lock1 should be holding the lock
			if !lock1.IsAcquired() {
				fmt.Printf("Lock1 failed to acquire lock\n\n")
			}
			lock1.Release()

			time.Sleep(time.Second * 3)

			return
		}
	}
}
