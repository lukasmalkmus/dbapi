package dbapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// A testRequest to test some functions.
type testRequest struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status int64  `json:"status"`
}

const (
	// testAPI is a test url that won`t be called.
	testAPI = "https://api.db.com/"

	// Test token
	testAccessToken = "1234567890abcdefghijklmnopqrstuvwxyz"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the dbapi client being tested.
	testClient *Client

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

func TestNew(t *testing.T) {
	api, err := New(
		SetToken(testAccessToken),
	)
	ok(t, err)

	// Is configuration present?
	equals(t, DefaultURL, api.baseURL.String())
	equals(t, DefaultVersion, api.version.String())
	equals(t, http.DefaultClient, api.client)
	equals(t, testAccessToken, api.Authentication.Token())

	// Are endpoints/resources present?
	equals(t, &AddressesService{client: api}, api.Addresses)
	equals(t, &AccountsService{client: api}, api.Accounts)
	equals(t, &TransactionsService{client: api}, api.Transactions)
	equals(t, &UserInfoService{client: api}, api.UserInfo)
}

func TestSetClient(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	mockData := []struct {
		HTTPClient         *http.Client
		ExpectedHTTPClient *http.Client
		ExpectedError      error
	}{
		{nil, http.DefaultClient, ErrInvalidClient},
		{customHTTPClient, customHTTPClient, nil},
	}

	for _, mock := range mockData {
		c, err := New(
			SetToken(testAccessToken),
			SetClient(mock.HTTPClient),
		)
		if c != nil {
			equals(t, c.client, mock.ExpectedHTTPClient)
		}
		equals(t, err, mock.ExpectedError)
	}
}

func TestSetToken(t *testing.T) {
	mockData := []struct {
		token         string
		ExpectedToken string
		ExpectedError error
	}{
		{"", "", nil},
		{"123", "123", nil},
		{testAccessToken, testAccessToken, nil},
	}

	for _, mock := range mockData {
		c, err := New(
			SetToken(testAccessToken),
			SetToken(mock.token),
		)
		if c != nil {
			equals(t, mock.ExpectedToken, c.Authentication.token)
		}
		equals(t, err, mock.ExpectedError)
	}
}

func TestSetURL(t *testing.T) {
	mockData := []struct {
		urlStr         string
		ExpectedURLStr string
		ExpectedError  error
	}{
		{testAPI, testAPI, nil},
		{"https://api.db.com", "https://api.db.com/", nil},
		{"", DefaultURL, ErrInvalidURL},
		{"://not-existing", DefaultURL, ErrInvalidURL},
	}

	for _, mock := range mockData {
		c, err := New(
			SetToken(testAccessToken),
			SetURL(mock.urlStr),
		)
		if c != nil {
			equals(t, mock.ExpectedURLStr, c.baseURL.String())
		}
		equals(t, err, mock.ExpectedError)
	}
}

func TestSetVersion(t *testing.T) {
	mockData := []struct {
		version            Version
		ExpectedVersion    Version
		ExpectedVersionStr string
		ExpectedError      error
	}{
		{DefaultVersion, V1, "v1", nil},
		{V1, V1, "v1", nil},
	}

	for _, mock := range mockData {
		c, err := New(
			SetToken(testAccessToken),
			SetVersion(mock.version),
		)
		if c != nil {
			equals(t, mock.ExpectedVersion, c.version)
			equals(t, mock.ExpectedVersionStr, c.version.String())
		}
		equals(t, err, mock.ExpectedError)
	}
}

func TestNewRequest(t *testing.T) {
	c, err := New(
		SetToken(testAccessToken),
		SetURL(testAPI),
	)
	ok(t, err)

	inURL, outURL := "/foo", testAPI+"v1/foo"
	inBody, outBody := &testRequest{ID: 1, Name: "Test Request", Status: 1}, `{"id":1,"name":"Test Request","status":1}`+"\n"
	req, _ := c.NewRequest("POST", inURL, inBody)

	// Test that relative URL was expanded.
	equals(t, outURL, req.URL.String())

	// Test that body was JSON encoded.
	body, _ := ioutil.ReadAll(req.Body)
	equals(t, outBody, string(body))
}

/*
func TestNewRequest_InvalidJSON(t *testing.T) {
		c, err := New(
			testAccessToken,
			SetURL(testAPI),
		)
		ok(t, err)

		type T struct {
			A map[int]interface{}
		}
		_, err = c.NewRequest(GET, "/", &T{})

		assert(t, err != nil, "Expected error to be returned.")

		if err, ok := err.(*json.UnsupportedTypeError); !ok {
			t.Errorf("Expected a JSON error; got %+v.", err)
		}
}
*/

func TestNewRequest_BadURL(t *testing.T) {

	c, err := New(
		SetToken(testAccessToken),
		SetURL(testAPI),
	)
	ok(t, err)

	_, err = c.NewRequest(GET, ":", nil)
	assert(t, err != nil, "Expected error to be returned.")

	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}

}

// If a nil body is passed to dbapi.NewRequest, make sure that nil is also
// passed to http.NewRequest. In most cases, passing an io.Reader that returns
// no content is fine, since there is no difference between an HTTP request body
// that is an empty string versus one that is not set at all. However in certain
// cases, intermediate systems may treat these differently resulting in subtile
// errors.
func TestNewRequest_EmptyBody(t *testing.T) {
	c, err := New(
		SetToken(testAccessToken),
		SetURL(testAPI),
	)
	ok(t, err)

	req, err := c.NewRequest(GET, "/", nil)
	ok(t, err)

	assert(t, req.Body == nil, "Constructed request contains a non-nil Body.")
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, GET)
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := testClient.NewRequest(GET, "/", nil)
	body := new(foo)
	testClient.Do(req, body)

	equals(t, &foo{"a"}, body)
}

func TestDo_ioWriter(t *testing.T) {
	setup()
	defer teardown()
	content := `{"A":"a"}`

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, GET)
		fmt.Fprint(w, content)
	})

	req, _ := testClient.NewRequest(GET, "/", nil)
	var buf []byte
	actual := bytes.NewBuffer(buf)
	testClient.Do(req, actual)

	equals(t, []byte(content), actual.Bytes())
}

func TestDo_HTTPError(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := testClient.NewRequest(GET, "/", nil)
	_, err := testClient.Do(req, nil)

	assert(t, err != nil, "Expected error to be returned (expected HTTP 400 error).")
}

// Test handling of an error caused by the internal http client's Do() function.
// A redirect loop is pretty unlikely to occur within the Cacheterrit API, but does allow us to exercise the right code path.
func TestDo_RedirectLoop(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := testClient.NewRequest(GET, "/", nil)
	_, err := testClient.Do(req, nil)

	assert(t, err != nil, "Expected error to be returned.")

	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

// setup sets up a test HTTP server along with a dbapi.Client that is configured
// to talk to that test server. Tests should register handlers on mux which
// provide mock responses for the API method being tested.
func setup() {
	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// dbapi client configured to use test server.
	testClient, _ = New(
		SetToken(testAccessToken),
		SetURL(testServer.URL),
	)
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	equals(t, want, r.Method)
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
