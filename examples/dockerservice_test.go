package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	".."
)

func setupFunc(stopper *inttest_utils.ServiceStopper) error {
	es, err := inttest_utils.NewDockerContainer("fabregas/test_service:1.0")
	if err != nil {
		return err
	}
	stopper.Add(es)

	es.SetEnv(map[string]string{"MYENV": "ok"})
	es.SetPortsBinding(map[string]string{"5555": "0.0.0.0:5747"})
	err = es.Start("my-test-service")
	if err != nil {
		return err
	}

	// waiting service started
	time.Sleep(5 * time.Second)
	return nil
}

func TestMain(m *testing.M) {
	os.Exit(inttest_utils.RunIntegrationTest(m, setupFunc))
}

func TestDockerService(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:5747/bar")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if string(body) != "Hello, world" {
		t.Fatal(string(body))
	}
}
