package inttest_utils

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func waitService(url string) int {
	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err == nil {
			//	fmt.Println(resp)
			return resp.StatusCode
		}
		time.Sleep(time.Second)
	}
	panic("service not started")
}

func TestDockerContainer(t *testing.T) {
	_, err := NewDockerContainer("fabregas/test_service:0.1")
	if err == nil {
		t.Fatal("expected error")
	}

	es, err := NewDockerContainer("fabregas/test_service:1.0")
	if err != nil {
		t.Fatal(err)
	}
	defer es.Stop()

	es.SetEnv(map[string]string{"MYENV": "ok"})
	es.SetPortsBinding(map[string]string{"5555": "0.0.0.0:5747"})
	err = es.Start("my-image-test-dc")
	if err != nil {
		t.Fatal(err)
	}

	if waitService("http://127.0.0.1:5747/bar") != 200 {
		t.Fatal("invalid resp status")
	}

	out, err := es.Logs()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "starting service..." {
		t.Fatal(out)
	}

	ip, err := es.IP()
	if err != nil {
		t.Fatal(err)
	}
	if len(ip) < 7 {
		t.Fatal("invalid IP addr")
	}

	err = es.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
