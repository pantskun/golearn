package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/pantskun/golearn/customEtcdclt/etcdinteraction"
)

func TestParseRequest(t *testing.T) {
	type httpTestInfo struct {
		query          string
		expectedAction etcdinteraction.EtcdActionInterface
	}

	r := httptest.NewRequest("GET", "http://localhost:8080/put", nil)
	putCases := []httpTestInfo{
		{"", etcdinteraction.NewPutAction("", "")},
		{"key=key", etcdinteraction.NewPutAction("key", "")},
		{"key=key&value=value", etcdinteraction.NewPutAction("key", "value")},
	}

	for _, putCase := range putCases {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(putCase.query))

		action := parsePutRequest(r)
		if !action.Equal(putCase.expectedAction) {
			t.Errorf("expect: %s, get: %s", putCase.expectedAction, action)
		}
	}

	r = httptest.NewRequest("GET", "http://localhost:8080/get", nil)
	getCases := []httpTestInfo{
		{"", etcdinteraction.NewGetAction("", "")},
		{"key=key", etcdinteraction.NewGetAction("key", "")},
	}

	for _, getCase := range getCases {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(getCase.query))

		action := parseGetRequest(r)
		if !action.Equal(getCase.expectedAction) {
			t.Errorf("expect: %s, get: %s", getCase.expectedAction, action)
		}
	}

	r = httptest.NewRequest("GET", "http://localhost:8080/del", nil)
	delCases := []httpTestInfo{
		{"", etcdinteraction.NewDeleteAction("", "")},
		{"key=key", etcdinteraction.NewDeleteAction("key", "")},
	}

	for _, delCase := range delCases {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(delCase.query))

		action := parseDeleteRequest(r)
		if !action.Equal(delCase.expectedAction) {
			t.Errorf("expect: %s, get: %s", delCase.expectedAction, action)
		}
	}
}

func TestWriteResponse(t *testing.T) {
	type writeTestInfo struct {
		msgs     []string
		err      error
		expected string
	}

	w := httptest.NewRecorder()

	writeCases := []writeTestInfo{
		{nil, etcdinteraction.EtcdError{Msg: "test"}, "test\n"},
		{[]string{"test1", "test2"}, nil, "test1\ntest2\n"},
	}

	for _, writeCase := range writeCases {
		writeResponse(writeCase.msgs, writeCase.err, w)

		body, _ := ioutil.ReadAll(w.Body)
		bodystr := string(body)

		if bodystr != writeCase.expected {
			t.Errorf("expect: %s, get: %s", writeCase.expected, bodystr)
		}
	}
}

func TestHTTPClient(t *testing.T) {
	wg := sync.WaitGroup{}

	var err error

	wg.Add(1)

	go func() {
		defer wg.Done()

		err = HTTPClient(":8080")
	}()

	time.Sleep(2 * time.Second)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	wg.Wait()

	if err.Error() != "http: Server closed" {
		t.Error("not got server close")
	}
}

// func TestListenSystemSignal(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go func(){

// 	}
// }
