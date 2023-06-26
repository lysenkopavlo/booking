package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	parameters         []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"cont", "/contacts", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"rm", "/reservation-summary", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	//creating a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, v := range theTest {
		if v.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + v.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != v.expectedStatusCode {
				t.Errorf("for %s expected status code is: %d, but we got: %d", v.name, v.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
