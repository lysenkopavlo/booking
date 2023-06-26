package main

import "testing"

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Errorf("failed run(): %v", err)
	}
}

//	to display test coverage in more details
//	use this command:
//	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
