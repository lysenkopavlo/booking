package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {

	case http.Handler:
		//Do nothing
	default:
		t.Errorf("wrong type of %T", v)
	}
}
func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {

	case http.Handler:
		//Do nothing
	default:
		t.Errorf("wrong type of %T", v)
	}
}

//	to display test coverage in more details
//	use this command:
//	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
