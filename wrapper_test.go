package inttest_utils

import (
	"fmt"
	"testing"
)

type RunnerTest struct{}

func (r *RunnerTest) Run() int {
	return 0
}

func TestWrapper(t *testing.T) {
	setupFunc := func(stopper *ServiceStopper) error {
		// starting some docker service
		es, err := NewDockerContainer("fabregas/test_service:1.0")
		if err != nil {
			return err
		}
		stopper.Add(es)

		es.SetEnv(map[string]string{"MYENV": "ok"})
		es.SetPortsBinding(map[string]string{"5555": "0.0.0.0:5747"})
		err = es.Start("my-image-test-dc")
		if err != nil {
			return err
		}

		// starting service as a regular application
		app, err := StartApplication("/tmp/test_service", map[string]string{"MYENV": "ok"}, false)
		if err != nil {
			return err
		}
		stopper.Add(app)

		if waitService("http://127.0.0.1:5747/bar") != 200 {
			return fmt.Errorf("docker service does not respond")
		}
		if waitService("http://127.0.0.1:5555/bar") != 200 {
			return fmt.Errorf("local service does not respond")
		}
		return nil
	}

	retcode := RunIntegrationTest(&RunnerTest{}, setupFunc)
	if retcode != 0 {
		t.Fatal(retcode)
	}
}
