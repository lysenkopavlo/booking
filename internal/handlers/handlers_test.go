package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lysenkopavlo/booking/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	// NOTE! For testing purpose we commented out from line 27 to line 48
	// Lesson No_120
	// And declare new func TestRepository_Reservation

	// tests for "GET" methods
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contacts", "GET", http.StatusOK},

	// tests for "post" methods
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

	for _, v := range theTests {
		resp, err := ts.Client().Get(ts.URL + v.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != v.expectedStatusCode {
			t.Errorf("for %s expected %d, but got %d", v.name, v.expectedStatusCode, resp.StatusCode)
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

	// Test case whean reservation is not into session
	// That's why we reseting variables
	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	ctx = getCTX(req) // now we have a contex which can be added to request

	req = req.WithContext(ctx) // and now we have request that knows about "X-Session" header

	// Simmulating request-response cycle down here
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	// But we don't put a reservation variable into context
	// coz we are testing redirection
	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nerror while testing session without reservation\nReservation handler returned wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test case whean roomID is wrong
	// That's why we reseting variables
	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	ctx = getCTX(req) // now we have a contex which can be added to request

	req = req.WithContext(ctx) // and now we have request that knows about "X-Session" header

	// Simmulating request-response cycle down here
	rr = httptest.NewRecorder()
	// But this time we put into context updated reservation variable with defenatly invalid roomID
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nerror while testing wrong roomID\nReservation handler returned wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler returned a wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing request body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler didn't parse anything\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid start_date
	// TODO: create a table test
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler didn't parse start_date\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end_date
	// TODO: create a table test
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler failed parsing end_date\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid room_di
	// TODO: create a table test
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler failed because of invalid room_id\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid first_name
	// TODO: create a table test
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-1-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler failed because of invalid name input\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for InsertReservation
	// TODO: create a table test
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-1-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler failed inserting a reservation\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for InsertRestriction
	// TODO: create a table test
	reqBody = "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-1-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Week")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=baba@yaga.ru")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=849627")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostReservation handler failed inserting a restriction\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := "start=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-03")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostAvailability handler returned a wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// without body ins session
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostAvailability handler returned a wrong response code\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// invalid start_date
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-03")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostAvailability handler cannot parse a start_date\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
	}

	// invalid start_date
	reqBody = "start=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("\nPostAvailability handler cannot parse a end_date\nGot: %d\twanted: %d", rr.Code, http.StatusSeeOther)
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

//	to display test coverage in more details
//	use this command:
//	go test -coverprofile=coverage.out && go tool cover -html=coverage.out
// 	if HTML doesn't open use `-func` flag
