package dbapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAddressesService_Get(t *testing.T) {
	setup()
	defer teardown()

	exp := &Addresses{
		{City: "Frankfurt", HouseNumber: 19, Street: "Große Bockenheimer Straße", Type: "MAILING_ADDRESS", ZipCode: 60311},
		{City: "Frankfurt", HouseNumber: 19, Street: "Große Bockenheimer Straße", Type: "REGISTRATION_ADDRESS", ZipCode: 60311},
	}

	testMux.HandleFunc("/v1/addresses", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"city":"Frankfurt","houseNumber":"19","street":"Große Bockenheimer Straße","type":"MAILING_ADDRESS","zip":"60311"},{"city":"Frankfurt","houseNumber":"19","street":"Große Bockenheimer Straße","type":"REGISTRATION_ADDRESS","zip":"60311"}]`)
	})

	act, _, err := testClient.Addresses.Get()
	ok(t, err)
	equals(t, exp, act)
}
