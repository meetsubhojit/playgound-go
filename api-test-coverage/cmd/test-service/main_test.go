//go:build api_test
// +build api_test

package main

import (
	"log"
	"net/http"
	"testing"

	"api-test-coverage/internal/server"
)

func TestRun(t *testing.T) {
	t.Log("starting test")
	go server.StartServer()

	t.Log("running tests")
	runApiTest()
	t.Log("test run done, exiting....")
}

// treat below function as running a API test using some framework like karate etc
func runApiTest() {
	resp, err := http.Get("http://localhost:8080/time")
	if err != nil {
		log.Fatal("got error in running test, err: " + err.Error())
	}
	if resp != nil {
		defer resp.Body.Close()
	}
}
