// To properly run test for render package
// We have to create session with context in it

package render

import (
	"net/http"
	"testing"

	"github.com/lysenkopavlo/booking/internal/models"
)

func TestAddData(t *testing.T) {
	td := models.TemplateData{}

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddData(&td, r)

	if result.Flash != "123" {
		t.Errorf("%v", result)
	}
}

func TestRenderTemplate(t *testing.T) {
	// This variable declared in the package scope in render.go
	// We assigned a new value here
	pathTemplates = "./../../templates"

	// Before testing a Render template func
	// we create a template cache

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Can't create template cache, error is: %v", err)
	}

	// If we've created template cache
	// we put it in app variable
	// which declared in the render.go
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var mw myWriter

	err = RenderTemplate(&mw, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Errorf("error writing template to browser, error: %v", err)
	}

	err = RenderTemplate(&mw, r, "non-existance.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Errorf("error writing non-existance template to browser, error: %v", err)
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-random-url", nil)

	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewTemplate(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
