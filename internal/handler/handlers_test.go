package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	// tests for "GET" methods
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"cont", "/contacts", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// test for "post" methods
	{"s-a", "/search-availability", "POST", []postData{
		{key: "start", value: "2022-06-26"},
		{key: "end", value: "2022-06-27"},
	}, http.StatusOK},
	{"s-a", "/search-availability", "POST", []postData{
		{key: "start", value: "2022-06-26"},
		{key: "end", value: "2022-06-27"},
	}, http.StatusOK},
	{"m-r", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "496-496-49688"},
	}, http.StatusOK},
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
		} else {
			values := url.Values{}
			for _, p := range v.parameters {
				values.Add(p.key, p.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+v.url, values)
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
