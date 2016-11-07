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
	// GET is a shortcut for http.MethodGet.
	GET = http.MethodGet
)

const (
	// DefaultURL is the URL of the Deutsche Bank API which is used by default.
	// NOTE: In a later release the version will be removed from the default URL
	// and needs to be passed explicitly.
	DefaultURL = "https://simulator-api.db.com/gw/dbapi/v1/"
)

var (
	ErrInvalidClient = errors.New("")
)

// A Client manages communication with the Deutsche Bank API.
type Client struct {
	client  *http.Client
	baseURL *url.URL
	token   string

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

// An Option serves as a 'functional parameter' which can be used to configure
// the behaviour of the API Client.
type Option func(c *Client) error

// SetClient specifies a custom http client that should be used to make requests.
func SetClient(client *http.Client) Option {
	return func(c *Client) error { return c.setClient(client) }
}
func (c *Client) setClient(client *http.Client) error {
	if client == nil {
		return errors.New("")
	}
	c.client = client
	return nil
}

// New creates and returns a new API Client. Options can be passed to configure
// the Client.
func New(AccessToken string, options ...Option) (*Client, error) {
	/*
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 3 * time.Second,
				TLSHandshakeTimeout:   3 * time.Second,
			},
		}
	*/

	// Parse the DefaultURL.
	url, err := url.Parse(DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create client with default settings.
	c := &Client{
		client:  http.DefaultClient,
		baseURL: url,
		token:   AccessToken,
	}
	c.Addresses = &AddressesService{client: c}
	c.Accounts = &AccountsService{client: c}
	c.Transactions = &TransactionsService{client: c}
	c.UserInfo = &UserInfoService{client: c}

	// Apply supplied options.
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
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

	// Apply Authentication (OAuth2 Access Token).
	// Documentation: https://developer.db.com/#/apidocumentation/apiauthorizationguide
	req.Header.Add("Authorization", "Bearer "+c.token)

	req.Header.Add("Accept", "application/json")

	return req, nil
}

// buildURLForRequest will build the URL (as string) that will be called. It
// does several cleaning tasks for us.
func (c *Client) buildURLForRequest(urlStr string) (string, error) {
	u := c.baseURL.String()

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
