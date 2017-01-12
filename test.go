//usr/bin/env go run $0 "$@"; exit
package main

import (
	"fmt"
	"os"
	"time"

	//apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/util/clock"

	"github.com/golang/glog"
)

type AutoScaler struct {
	pollPeriod    time.Duration
	clock         clock.Clock
	stopCh        chan struct{}
	readyCh       chan<- struct{} // For testing.
}

func NewAutoScaler() (*AutoScaler, error) {
	return &AutoScaler{
		pollPeriod:    time.Second * time.Duration(5),
		clock:         clock.RealClock{},
		stopCh:        make(chan struct{}),
		readyCh:       make(chan struct{}, 1),
	}, nil
}

func main() {
	scaler, err := NewAutoScaler()
	if err != nil {
		glog.Errorf("%v", err)
		os.Exit(1)
	}
	// Begin autoscaling.
	scaler.Run()
}

func (s *AutoScaler) Run() {
	ticker := s.clock.Tick(s.pollPeriod)
	s.readyCh <- struct{}{} // For testing.

	// Don't wait for ticker and execute pollPrometheus() for the first time.
	s.pollPrometheus()

	for {
		select {
		case <-ticker:
			s.pollPrometheus()
		case <-s.stopCh:
			return
		}
	}
}

func (s *AutoScaler) pollPrometheus() {
	fmt.Printf("\n%v\n", time.Now().Unix())
}
