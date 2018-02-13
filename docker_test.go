package inttest_utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	// create tmp file for mounting into container
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("/tmp", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// gettting test image
	_, err = NewDockerContainer("fabregas/test_service:0.1")
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
	es.SetVolumesBinding(map[string]string{tmpfile.Name(): "/container_volume"})
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

	// check mount
	resp, err := http.Get("http://127.0.0.1:5747/volume/size")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != fmt.Sprintf("%d", len(content)) {
		t.Fatal(string(body))
	}

	err = es.Stop()
	if err != nil {
		t.Fatal(err)
	}

}
