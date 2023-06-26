package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/lysenkopavlo/booking/internal/config"
	"github.com/lysenkopavlo/booking/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am igoing to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	// Session is a variable to set user's session parameters
	session = scs.New()

	// setting lifetime for session
	session.Lifetime = 24 * time.Hour

	// is cockie persisting after closing a browser?
	session.Cookie.Persist = true // means don't want to clear cockie after closing browser

	// how restrict do you want to be about coockie?
	session.Cookie.SameSite = http.SameSiteLaxMode

	// insisting about coockie incryption
	session.Cookie.Secure = false // for now

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

// Creating a type myWriter
// will mimic the http.ResponceWriter interface
// By implementing it's methods
type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	return http.Header{}
}

func (mw *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw *myWriter) WriteHeader(statusCode int) {

}
