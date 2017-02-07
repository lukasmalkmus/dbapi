package dbapi

import (
	"fmt"
	"net/http"
)

// The TransactionsService binds to the HTTP endpoints which belong to
// the transactions resource.
type TransactionsService struct {
	client *Client
}

// Transactions are the users transactions.
type Transactions []struct {
	Amount           float64 `json:"amount,omitempty"`
	CounterPartyName string  `json:"counterPartyName,omitempty"`
	CounterPartyIBAN string  `json:"counterPartyIban,omitempty"`
	Usage            string  `json:"usage,omitempty"`
	Date             string  `json:"date,omitempty"`
}

// GetAll reads all transactions of all accounts of the current user. It is
// not apparent who issued a transaction, only whether the user gained or lost
// money by it (based on wether the amount is positive or negative respectively).
func (s *TransactionsService) GetAll() (*Transactions, *Response, error) {
	u := "/transactions"
	r := new(Transactions)

	resp, err := s.client.Call(http.MethodGet, u, nil, r)
	return r, resp, err
}

// Get all transactions for a specific account of the current user. If given
// IBAN is not valid or does not represent an account of the current user, an
// empty result is returned. It is not apparent who issued a transaction, only
// whether the user gained or lost money by it (based on wether the amount is
// positive or negative respectively).
func (s *TransactionsService) Get(iban string) (*Transactions, *Response, error) {
	u := fmt.Sprintf("/transactions?iban=%s", iban)
	r := new(Transactions)

	resp, err := s.client.Call(http.MethodGet, u, nil, r)
	return r, resp, err
}
