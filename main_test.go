package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RaceCondition(t *testing.T) {

	expected := []string{"Knock knock\n", "Who's there\n", "Race condtion\n"}
	results := []string{}

	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	callServer := func(path string) {
		result, err := responseStringFromServerCall(path)
		assert.Nil(t, err)
		results = append(results, result)
	}

	go callServer(srv.URL + "/line1")
	go callServer(srv.URL + "/line2")
	go callServer(srv.URL + "/line3")

	time.Sleep(time.Second)

	assert.Equal(t, expected, results)
}

// func Test_RaceCondition(t *testing.T) {

// 	expected := []string{"Knock knock\n", "Who's there\n", "Race condtion\n"}
// 	results := []string{}

// 	addr, err := getFreePort()
// 	assert.Nil(t, err)
// 	srv := StartServer(addr)
// 	defer srv.Close()

// 	callServer := func(path string) {
// 		result, err := responseStringFromServerCall(path)
// 		assert.Nil(t, err)
// 		results = append(results, result)
// 	}

// 	go callServer(addr + "/line1")
// 	go callServer(addr + "/line2")
// 	go callServer(addr + "/line3")

// 	time.Sleep(time.Second)

// 	assert.Equal(t, expected, results)
// }

func Test_EndsBeforeFailing(t *testing.T) {

	go func() {
		assert.Nil(t, errors.New("uh oh"))
	}()
}

func Test_RaceCondition_fixed(t *testing.T) {

	expected := []string{"Knock knock\n", "Who's there\n", "Race condtion\n"}
	var resultsMu sync.Mutex
	results := []string{}

	srv := httptest.NewServer(GetRouter())
	defer srv.Close()

	var wg sync.WaitGroup
	wg.Add(len(expected))

	errCh := make(chan error)

	callServer := func(path string) {
		result, err := responseStringFromServerCall(path)
		if err != nil {
			errCh <- err
			return
		}
		resultsMu.Lock()
		results = append(results, result)
		resultsMu.Unlock()
		wg.Done()
		time.Sleep(time.Millisecond)
	}

	go callServer(srv.URL + "/line1")
	go callServer(srv.URL + "/line2")
	go callServer(srv.URL + "/line3")

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		t.Fatalf("got error %v", err)
	case <-time.After(time.Second):
		t.Fatal("waited too long for http requests")
	case <-done:
	}

	resultsMu.Lock()
	defer resultsMu.Unlock()
	assert.Len(t, results, len(expected))
	for _, sentence := range expected {
		assert.Contains(t, results, sentence)
	}

}

// Test Helpers

func responseStringFromServerCall(path string) (string, error) {
	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, 1024)
	n, err := resp.Body.Read(bytes)
	if err != io.EOF {
		return "", err
	}
	return string(bytes[0:n]), nil
}
