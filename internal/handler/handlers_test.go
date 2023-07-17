package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/lysenkopavlo/booking/internal/models"
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
	// NOTE! For testing purpose we commented out from line 27 to line 48
	// Lesson No_120
	// And declare new func TestRepository_Reservation

	// tests for "GET" methods
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	// {"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	// {"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"cont", "/contacts", "GET", []postData{}, http.StatusOK},
	// {"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// // tests for "post" methods
	// {"s-a", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2022-06-26"},
	// 	{key: "end", value: "2022-06-27"},
	// }, http.StatusOK},
	// {"s-a", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2022-06-26"},
	// 	{key: "end", value: "2022-06-27"},
	// }, http.StatusOK},
	// {"m-r", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "496-496-49688"},
	// }, http.StatusOK},
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

func TestRepository_Reservation(t *testing.T) {
	// Because func Reservation from handler package uses reservation data
	// Of type models.Reservation, which exracted from the session,
	// We create a dummy version of this
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	// In app after extraction reservation data from session
	// We check roomName by its ID
	// And then put updated reservation into session back
	// We simmulate this proccess here:
	req, _ := http.NewRequest("GET", "/make-reservation", nil)

	ctx := getCTX(req) // now we have a contex which can be added to request

	req = req.WithContext(ctx) // and now we have request that knows about "X-Session" header

	// Simmulating request-response cycle down here
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("\nReservation handler returned wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusOK)
	}
}

// getCTX returns contex for TestRepository_Reservation purposes
func getCTX(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session")) // It's neccessary to write "X-Session" coz of testing
	if err != nil {
		log.Println(err)
	}
	return ctx
}
