package dbapi

import "net/http"

// The AddressesService binds to the HTTP endpoints which belong to the
// addresses resource.
type AddressesService struct {
	client *Client
}

// Addresses are the users addresses.
type Addresses []struct {
	Street      string `json:"street,omitempty"`
	HouseNumber int64  `json:"houseNumber,string,omitempty"`
	ZipCode     int64  `json:"zip,string,omitempty"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	Type        string `json:"type,omitempty"`
}

// Get reads all addresses of the current user. Usually a user has exactly two
// addresses with the types MAILING_ADDRESS and REGISTRATION_ADDRESS
// respectively. Otherwise those two addresses are often identical.
func (s *AddressesService) Get() (*Addresses, *Response, error) {
	u := "/addresses"
	r := new(Addresses)

	resp, err := s.client.Call(http.MethodGet, u, nil, r)
	return r, resp, err
}
