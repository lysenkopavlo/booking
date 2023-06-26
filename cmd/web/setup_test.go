package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHandler struct{}

func (myH myHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {

}

//	to display test coverage in more details
//	use this command:
//	go test -coverprofile=c.out && go tool cover -html=c.out
