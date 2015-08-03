package objdbClient

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/contiv/objmodel/objdb"

	log "github.com/Sirupsen/logrus"
)

type JsonObj struct {
	Value string
}

// New objdb client
var client = NewClient()

// Perform Set/Get operation on default conf store
func TestSetGet(t *testing.T) {
	runtime.GOMAXPROCS(4)

	// Set
	setVal := JsonObj{
		Value: "test1",
	}
	err := client.SetObj("/contiv.io/test", setVal)
	if err != nil {
		fmt.Printf("Error setting key. Err: %v\n", err)
		t.Errorf("Error setting key")
	}

	var retVal JsonObj
	err = client.GetObj("/contiv.io/test", &retVal)
	if err != nil {
		fmt.Printf("Error getting key. Err: %v\n", err)
		t.Errorf("Error getting key")
	}

	if retVal.Value != "test1" {
		fmt.Printf("Got invalid response: %+v\n", retVal)
		t.Errorf("Got invalid response")
	}

	err = client.DelObj("/contiv.io/test")
	if err != nil {
		t.Errorf("Error deleting test object. Err: %v", err)
	}

	fmt.Printf("Set/Get/Del test successful\n")
}

func TestLockAcquireRelease(t *testing.T) {
	// Create a lock
	lock1, err := client.NewLock("master", "hostname1", 10)
	lock2, err := client.NewLock("master", "hostname2", 10)

	// Acquire the lock
	err = lock1.Acquire(0)
	if err != nil {
		t.Errorf("Error acquiring lock1")
	}

	// Try to acquire the same lock again. This should fail
	err = lock2.Acquire(0)
	if err != nil {
		t.Errorf("Error acquiring lock2")
	}

	cnt := 1
	for {
		select {
		case event := <-lock1.EventChan():
			fmt.Printf("Event on Lock1: %+v\n\n", event)
			if event.EventType == objdb.LockAcquired {
				fmt.Printf("Master lock acquired by Lock1\n")
			}
		case event := <-lock2.EventChan():
			fmt.Printf("Event on Lock2: %+v\n\n", event)
			if event.EventType == objdb.LockAcquired {
				fmt.Printf("Master lock acquired by Lock2\n")
			}
		case <-time.After(time.Second * time.Duration(30)):
			if cnt == 1 {
				fmt.Printf("30sec timer. releasing Lock1\n\n")
				// At this point, lock1 should be holding the lock
				if !lock1.IsAcquired() {
					t.Errorf("Lock1 failed to acquire lock\n\n")
				}

				// Release lock1 so that lock2 can acquire it
				lock1.Release()
				cnt++
			} else {
				fmt.Printf("60sec timer. checking if lock2 is acquired\n\n")

				// At this point, lock2 should be holding the lock
				if !lock2.IsAcquired() {
					t.Errorf("Lock2 failed to acquire lock\n\n")
				}

				fmt.Printf("Success. Lock2 Successfully acquired. releasing it\n")
				// we are done with the test
				lock2.Release()

				return
			}
		}
	}
}

func TestLockAcquireTimeout(t *testing.T) {
	fmt.Printf("\n\n\n =========================================================== \n\n\n")
	// Create a lock
	lock1, err := client.NewLock("master", "hostname1", 10)
	lock2, err := client.NewLock("master", "hostname2", 10)

	// Acquire the lock
	err = lock1.Acquire(0)
	if err != nil {
		t.Errorf("Error acquiring lock1")
	}

	time.Sleep(300 * time.Millisecond)

	err = lock2.Acquire(20)
	if err != nil {
		t.Errorf("Error acquiring lock2")
	}

	for {
		select {
		case event := <-lock1.EventChan():
			fmt.Printf("Event on Lock1: %+v\n\n", event)
			if event.EventType == objdb.LockAcquired {
				fmt.Printf("Master lock acquired by Lock1\n")
			}
		case event := <-lock2.EventChan():
			fmt.Printf("Event on Lock2: %+v\n\n", event)
			if event.EventType != objdb.LockAcquireTimeout {
				fmt.Printf("Invalid event on Lock2\n")
			} else {
				fmt.Printf("Lock2 timeout as expected")
			}
		case <-time.After(time.Second * time.Duration(40)):
			fmt.Printf("40sec timer. releasing Lock1\n\n")
			// At this point, lock1 should be holding the lock
			if !lock1.IsAcquired() {
				t.Errorf("Lock1 failed to acquire lock\n\n")
			}
			lock1.Release()

			time.Sleep(time.Second * 3)

			return
		}
	}
}

func TestServiceRegister(t *testing.T) {
	// Service info
	service1Info := objdb.ServiceInfo{"athena", "10.10.10.10", 4567}
	service2Info := objdb.ServiceInfo{"athena", "10.10.10.10", 4568}

	// register it
	err := client.RegisterService(service1Info)
	if err != nil {
		t.Errorf("Error registering service. Err: %+v\n", err)
	}
	log.Infof("Registered service: %+v", service1Info)

	err = client.RegisterService(service2Info)
	if err != nil {
		t.Errorf("Error registering service. Err: %+v\n", err)
	}
	log.Infof("Registered service: %+v", service2Info)

	resp, err := client.GetService("athena")
	if err != nil {
		t.Errorf("Error getting service. Err: %+v\n", err)
	}

	log.Infof("Got service list: %+v\n", resp)

	if (len(resp) < 2) || (resp[0] != service1Info) || (resp[1] != service2Info) {
		t.Errorf("Resp service list did not match input")
	}

	// Wait a while to make sure background refresh is working correctly
	time.Sleep(time.Second * 90)

	resp, err = client.GetService("athena")
	if err != nil {
		t.Errorf("Error getting service. Err: %+v\n", err)
	}

	log.Infof("Got service list: %+v\n", resp)

	if (len(resp) < 2) || (resp[0] != service1Info) || (resp[1] != service2Info) {
		t.Errorf("Resp service list did not match input")
	}
}

func TestServiceDeregister(t *testing.T) {
	// Service info
	service1Info := objdb.ServiceInfo{"athena", "10.10.10.10", 4567}
	service2Info := objdb.ServiceInfo{"athena", "10.10.10.10", 4568}

	// register it
	err := client.DeregisterService(service1Info)
	if err != nil {
		t.Errorf("Error deregistering service. Err: %+v\n", err)
	}
	err = client.DeregisterService(service2Info)
	if err != nil {
		t.Errorf("Error deregistering service. Err: %+v\n", err)
	}

	time.Sleep(time.Second * 10)
}

func TestServiceWatch(t *testing.T) {
	service1Info := objdb.ServiceInfo{"athena", "10.10.10.10", 4567}

	// register it
	err := client.RegisterService(service1Info)
	if err != nil {
		t.Errorf("Error registering service. Err: %+v\n", err)
	}
	log.Infof("Registered service: %+v", service1Info)

	// Create event channel
	eventChan := make(chan objdb.WatchServiceEvent, 1)
	stopChan := make(chan bool, 1)

	// Start watching for service
	err = client.WatchService("athena", eventChan, stopChan)
	if err != nil {
		t.Errorf("Error watching service. Err %v", err)
	}

	cnt := 1
	for {
		select {
		case srvEvent := <-eventChan:
			log.Infof("\n----\nReceived event: %+v\n----", srvEvent)
		case <-time.After(time.Second * time.Duration(10)):
			service2Info := objdb.ServiceInfo{"athena", "10.10.10.11", 4567}
			if cnt == 1 {
				// register it
				err := client.RegisterService(service2Info)
				if err != nil {
					t.Errorf("Error registering service. Err: %+v\n", err)
				}
				log.Infof("Registered service: %+v", service2Info)
			} else if cnt == 5 {
				// deregister it
				err := client.DeregisterService(service2Info)
				if err != nil {
					t.Errorf("Error deregistering service. Err: %+v\n", err)
				}
				log.Infof("Dregistered service: %+v", service2Info)
			} else if cnt == 7 {
				// Stop the watch
				stopChan <- true

				// wait a little and exit
				time.Sleep(time.Second)

				return
			}
			cnt++
		}
	}
}
