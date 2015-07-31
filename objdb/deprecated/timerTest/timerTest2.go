package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

var counter uint64 = 0

type Obj struct {
	Val       uint64
	EventChan chan bool
	WaitChan  chan bool
}

func waitLoop(obj *Obj) {
	time.Sleep(time.Second * 3)
	for {
		select {
		case <-obj.WaitChan:
			log.Infof("Received on wait Loop for %d\n", obj.Val)
		case <-time.After(time.Second * time.Duration(7)):
			log.Infof("Timeout on wait Loop for %d\n", obj.Val)
			obj.EventChan <- true
		}
	}
}

func runLoop(obj *Obj) {
	if (obj.Val % 2) == 0 {
		waitLoop(obj)
	}

	for {
		select {
		case <-obj.WaitChan:
			log.Infof("Received on wait chan for %d\n", obj.Val)
		case <-time.After(time.Second * time.Duration(7)):
			log.Infof("Timeout on wait chan for %d\n", obj.Val)
		}
	}
}

func NewLoop() *Obj {
	obj := new(Obj)

	obj.Val = counter
	counter = counter + 1
	obj.EventChan = make(chan bool, 1)
	obj.WaitChan = make(chan bool, 1)

	go runLoop(obj)

	return obj
}

func init() {
	log.Infof("Running timer test\n")

	obj1 := NewLoop()
	obj2 := NewLoop()

	for {
		select {
		case <-obj1.EventChan:
			log.Infof("Received event on Obj1\n")
		case <-obj2.EventChan:
			log.Infof("Received event on Obj2\n")
		case <-time.After(time.Second * time.Duration(5)):
			log.Infof("Received timer Event\n")
		}
	}
}
