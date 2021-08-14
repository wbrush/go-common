package helpers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/urfave/negroni"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
)

type (
	HandlerFunc   func(http.ResponseWriter, *http.Request)
	BodyCheckFunc func(*bytes.Buffer) error
)

type HandlerTestData struct {
	T             *testing.T
	Method        string
	URLMask       string
	URL           string
	Headers       map[string]string
	Body          io.Reader
	Query         interface{}
	Middleware    []negroni.HandlerFunc
	HandlerFunc   HandlerFunc
	ExpStatus     int
	BodyCheckFunc BodyCheckFunc
}

func HelperSchemaEncodeMust(src interface{}) url.Values {
	dst := make(map[string][]string)
	_ = schema.NewEncoder().Encode(src, dst)
	return url.Values(dst)
}

func HelperJsonMarshallMust(v interface{}) string {
	d, _ := json.Marshal(v)
	return string(d)
}

func HelperHttpHandlerTest(d HandlerTestData) {
	// Create a request to pass to our handler
	r, err := http.NewRequest(d.Method, d.URL, d.Body)
	if err != nil {
		d.T.Fatal(err)
	}

	if d.Query != nil {
		v, ok := d.Query.(url.Values)
		if !ok {
			d.T.Fatalf("bad query values provided")
		}
		r.URL.RawQuery = v.Encode()
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	w := httptest.NewRecorder()

	// check to see if we can get json structure back as well
	r.Header.Add("Accept", "application/json")

	//add additional headers if provided
	if len(d.Headers) > 0 {
		for k, v := range d.Headers {
			r.Header.Add(k, v)
		}
	}

	if d.URLMask == "" { //set URL mask to URL if not provided
		d.URLMask = d.URL
	}

	// we create mux router to process our requests
	// it's need for correctly process URL-params
	router := mux.NewRouter()

	wrapper := negroni.New()
	with := wrapper.With()
	for _, m := range d.Middleware {
		with.Use(m)
	}

	with.Use(negroni.Wrap(http.HandlerFunc(d.HandlerFunc)))
	router.Handle(d.URLMask, with)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(w, r)

	// Check the status code is what we expect.
	if status := w.Code; status != d.ExpStatus {
		d.T.Errorf("handler returned wrong status code: got %v want %v",
			status, d.ExpStatus)
	}

	// Check the response body is what we expect.
	if err := d.BodyCheckFunc(w.Body); err != nil {
		d.T.Errorf("handler returned unexpected body: %s", err.Error())
	}
}
