package dbapi

import (
	"fmt"
	"net/http"
)

// The AccountsService binds to the HTTP endpoints which belong to the
// cashAccounts resource.
type AccountsService struct {
	client *Client
}

// Accounts are the cash accounts of the user.
type Accounts []struct {
	Iban               string  `json:"iban,omitempty"`
	Balance            float64 `json:"balance,omitempty"`
	ProductDescription string  `json:"productDescription,omitempty"`
}

// GetAll reads all cash accounts of the current user. Only current accounts and
// accounts in the currency EUR are returned.
func (s *AccountsService) GetAll() (*Accounts, *Response, error) {
	u := "/cashAccounts"
	r := new(Accounts)

	resp, err := s.client.Call(http.MethodGet, u, nil, r)
	return r, resp, err
}

// Get reads the specified cash account of the current user. If given IBAN is
// not valid or does not represent an account of the current user, an empty
// result is returned.
func (s *AccountsService) Get(iban string) (*Accounts, *Response, error) {
	u := fmt.Sprintf("/cashAccounts?iban=%s", iban)
	r := new(Accounts)

	resp, err := s.client.Call(http.MethodGet, u, nil, r)
	return r, resp, err
}
