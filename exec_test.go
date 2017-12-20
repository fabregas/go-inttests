package inttest_utils

import (
	"os/exec"
	"testing"
	"time"
)

func TestExecApp(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "/tmp/test_service", "./test_service/main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	// start without env
	app, err := StartApplication("/tmp/test_service", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = app.Stop()
	if err != nil {
		t.Fatal(err)
	}
	if app.ExitCode != 1 {
		t.Fatal(app.ExitCode)
	}

	// start with env
	app, err = StartApplication("/tmp/test_service -f someflag", map[string]string{"MYENV": "ok"}, true)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Stop()

	if waitService("http://127.0.0.1:5555/bar") != 200 {
		t.Fatal("invalid resp status")
	}

	err = app.Stop()
	if err != nil {
		t.Fatal(err)
	}
	if app.ExitCode != -1 {
		t.Fatal(app.ExitCode)
	}
}
