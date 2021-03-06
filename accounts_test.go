package dbapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAccountsService_GetAll(t *testing.T) {
	setup()
	defer teardown()

	exp := &Accounts{
		{Iban: "DE10000000000000000453", Balance: 31236.95, ProductDescription: "persönliches Konto"},
		{Iban: "DE10000000000000000454", Balance: 250, ProductDescription: "persönliches Konto"},
		{Iban: "DE10000000000000000455", Balance: 100, ProductDescription: "persönliches Konto"},
	}

	testMux.HandleFunc("/v1/cashAccounts", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"iban":"DE10000000000000000453","balance":31236.95,"productDescription":"persönliches Konto"},{"iban":"DE10000000000000000454","balance":250,"productDescription":"persönliches Konto"},{"iban":"DE10000000000000000455","balance":100,"productDescription":"persönliches Konto"}]`)
	})

	act, _, err := testClient.Accounts.GetAll()
	ok(t, err)
	equals(t, exp, act)
}

func TestAccountsService_Get(t *testing.T) {
	setup()
	defer teardown()

	exp := &Accounts{
		{Iban: "DE10000000000000000454", Balance: 250, ProductDescription: "persönliches Konto"},
	}

	testMux.HandleFunc("/v1/cashAccounts", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"iban":"DE10000000000000000454","balance":250,"productDescription":"persönliches Konto"}]`)
	})

	act, _, err := testClient.Accounts.Get("DE10000000000000000454")
	ok(t, err)
	equals(t, exp, act)
}
