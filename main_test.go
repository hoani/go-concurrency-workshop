package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

func Test_EndsBeforeFailing(t *testing.T) {

	go func() {
		assert.Nil(t, errors.New("uh oh"))
	}()
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
