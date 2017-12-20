package inttest_utils

import (
	"fmt"
)

type Stopper interface {
	Stop() error
}

type Runner interface {
	Run() int
}

type ServiceStopper struct {
	services []Stopper
}

func (ss *ServiceStopper) Add(service Stopper) {
	ss.services = append(ss.services, service)
}
func (ss *ServiceStopper) StopAll() error {
	isErr := false
	for _, s := range ss.services {
		err := s.Stop()
		if err != nil {
			isErr = true
			fmt.Printf("[DESTROY ERROR] %s\n", err.Error())
		}
	}
	if isErr {
		return fmt.Errorf("stop services failed")
	}
	return nil
}

type SetupEnvFunc func(stopper *ServiceStopper) error

func RunIntegrationTest(r Runner, sef SetupEnvFunc) int {
	stopper := &ServiceStopper{make([]Stopper, 0)}
	defer stopper.StopAll()

	err := sef(stopper)
	if err != nil {
		fmt.Printf("Setup environment failed: %s\n", err.Error())
		return 1
	}

	return r.Run()
}
