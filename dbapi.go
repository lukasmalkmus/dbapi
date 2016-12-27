package dbapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// Version is the version of this package.
	version   = "0.4.2"
	userAgent = "dbapi/" + version
)

const (
	// GET is a shortcut for http.MethodGet.
	GET = http.MethodGet
)

const (
	// V1 is the version 1 of the dbAPI.
	V1 = "v1"
)

const (
	// DefaultURL is the URL of the Deutsche Bank API which is used by default.
	DefaultURL = "https://simulator-api.db.com/gw/dbapi/"

	// DefaultVersion is the default API version to use and defaults to v1.
	DefaultVersion = V1
)

var (
	// ErrInvalidClient is raised when the costum HTTP client is invalid (e.g.
	// nil).
	ErrInvalidClient = errors.New("Invalid http client!")

	// ErrInvalidURL is raised when the url couldn't be parsed by url.Parse().
	ErrInvalidURL = errors.New("Invalid url!")
)

// A Client manages communication with the Deutsche Bank API.
type Client struct {
	client  *http.Client
	baseURL *url.URL
	version Version

	// Authentication
	Authentication *AuthenticationService

	// API Resources
	Addresses    *AddressesService
	Accounts     *AccountsService
	Transactions *TransactionsService
	UserInfo     *UserInfoService
}

// A Response represents a http response from the Deutsche Bank API. It is a
// wrapper around the standard http.Response type.
type Response struct {
	*http.Response
}

// Version is the API version.
type Version string

func (v Version) String() string {
	return string(v)
}

// An Option serves as a 'functional parameter' which can be used to configure
// the behaviour of the API Client.
type Option func(c *Client) error

// SetClient specifies a custom http client that should be used to make requests.
// An error ErrInvalidClient is returned if the passed client is nil.
func SetClient(client *http.Client) Option {
	return func(c *Client) error { return c.setClient(client) }
}
func (c *Client) setClient(client *http.Client) error {
	if client == nil {
		return ErrInvalidClient
	}
	c.client = client
	return nil
}

// SetToken specifies the api token.
func SetToken(token string) Option {
	return func(c *Client) error { return c.setToken(token) }
}
func (c *Client) setToken(token string) error {
	c.Authentication.token = token
	return nil
}

// SetURL specifies the base url to use. An error ErrInvalidURL is returned if
// the passed url string can't be parsed properly.
func SetURL(urlStr string) Option {
	return func(c *Client) error { return c.setURL(urlStr) }
}
func (c *Client) setURL(urlStr string) error {
	if len(urlStr) == 0 {
		return ErrInvalidURL
	}
	// If there is no / at the end, add one.
	if strings.HasSuffix(urlStr, "/") == false {
		urlStr += "/"
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return ErrInvalidURL
	}
	c.baseURL = url
	return nil
}

// SetVersion specifies the api version to use.
func SetVersion(version Version) Option {
	return func(c *Client) error { return c.setVersion(version) }
}
func (c *Client) setVersion(version Version) error {
	c.version = version
	return nil
}

// New creates and returns a new API Client. Options can be passed to configure
// the Client.
func New(options ...Option) (*Client, error) {
	// Parse the DefaultURL.
	url, err := url.Parse(DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create client with default settings.
	c := &Client{
		client:  http.DefaultClient,
		baseURL: url,
		version: DefaultVersion,
	}
	c.Authentication = &AuthenticationService{}
	c.Addresses = &AddressesService{client: c}
	c.Accounts = &AccountsService{client: c}
	c.Transactions = &TransactionsService{client: c}
	c.UserInfo = &UserInfoService{client: c}

	// Apply supplied options.
	if err := c.Options(options...); err != nil {
		return nil, err
	}

	return c, nil
}

// Options applies Options to a client instance.
func (c *Client) Options(options ...Option) error {
	for _, option := range options {
		if err := option(c); err != nil {
			return err
		}
	}
	return nil
}

// Call combines Client.NewRequest() and Client.Do() methodes to avoid code
// duplication.
//
// m is the HTTP method you want to call.
// u is the URL you want to call.
// b is the HTTP body.
// r is the HTTP response.
//
// For more information read https://github.com/google/go-github/issues/234
func (c *Client) Call(m, u string, b interface{}, r interface{}) (*Response, error) {
	req, err := c.NewRequest(m, u, b)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, r)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	return fmt.Errorf("API call to %s failed: %s", r.Request.URL.String(), r.Status)
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by r, or
// returned as an error if an API error has occurred. If r implements the
// io.Writer interface, the raw response body will be written to r, without
// attempting to first decode it.
func (c *Client) Do(req *http.Request, r interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Wrap response
	response := &Response{Response: resp}

	err = CheckResponse(resp)
	if err != nil {
		// Return respone in case the caller wants to inspect it further.
		return response, err
	}

	if r != nil {
		if w, ok := r.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				// Return respone in case the caller wants to inspect it further.
				return response, err
			}
			err = json.Unmarshal(body, r)
		}
	}
	return response, err
}

// NewRequest creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved
// relative to the baseURL of the Client. Relative URLs should always be
// specified without a preceding slash. If specified, the value pointed to by
// body is JSON encoded and included as the request body.
func (c *Client) NewRequest(m, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.buildURLForRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(m, u, buf)
	if err != nil {
		return nil, err
	}

	// Apply Authentication if credentials are present.
	// Documentation: https://developer.db.com/#/apidocumentation/apiauthorizationguide
	if c.Authentication.HasAuth() {
		req.Header.Add("Authorization", "Bearer "+c.Authentication.Token())
	}

	// Add some important headers.
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

// buildURLForRequest will build the URL (as string) that will be called. It
// does several cleaning tasks for us.
func (c *Client) buildURLForRequest(urlStr string) (string, error) {
	u := c.baseURL.String() + c.version.String()

	// If there is no / at the end, add one.
	if strings.HasSuffix(u, "/") == false {
		u += "/"
	}

	// If there is a "/" at the start, remove it.
	if strings.HasPrefix(urlStr, "/") == true {
		urlStr = urlStr[1:]
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	u += rel.String()

	return u, nil
}
