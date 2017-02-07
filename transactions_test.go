package dbapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTransactionsService_GetAll(t *testing.T) {
	setup()
	defer teardown()

	exp := &Transactions{
		{Amount: -35.56, CounterPartyName: "Netto", Usage: "POS MIT PIN. Einkauf", Date: "2016-10-27"},
		{Amount: -52.22, CounterPartyName: "Lidl", Usage: "POS MIT PIN. Einkauf", Date: "2016-10-24"},
		{Amount: -1500, CounterPartyName: "Schwäbisch Hall", Usage: "Ref. 58974-8765889", Date: "2016-10-21"},
		{Amount: -38.98, CounterPartyName: "Toys R Us", Usage: "Rechnung", Date: "2016-10-17"},
		{Amount: -25.95, CounterPartyName: "Alnatura Frankfurt", Usage: "POS MIT PIN. Einkauf", Date: "2016-10-17"},
		{Amount: -96.16, CounterPartyName: "JET", Usage: "POS MIT PIN. Die Tanke Ihrer Wahl", Date: "2016-10-12"},
	}

	testMux.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"amount":-35.56,"counterPartyName":"Netto","usage":"POS MIT PIN. Einkauf","date":"2016-10-27"},{"amount":-52.22,"counterPartyName":"Lidl","usage":"POS MIT PIN. Einkauf","date":"2016-10-24"},{"amount":-1500,"counterPartyName":"Schwäbisch Hall","usage":"Ref. 58974-8765889","date":"2016-10-21"},{"amount":-38.98,"counterPartyName":"Toys R Us","usage":"Rechnung","date":"2016-10-17"},{"amount":-25.95,"counterPartyName":"Alnatura Frankfurt","usage":"POS MIT PIN. Einkauf","date":"2016-10-17"},{"amount":-96.16,"counterPartyName":"JET","usage":"POS MIT PIN. Die Tanke Ihrer Wahl","date":"2016-10-12"}]`)
	})

	act, _, err := testClient.Transactions.GetAll()
	ok(t, err)
	equals(t, exp, act)
}

func TestTransactionsService_Get(t *testing.T) {
	setup()
	defer teardown()

	exp := &Transactions{
		{Amount: 50, CounterPartyName: "Claudia Klar", Usage: "Sparen Samuel", Date: "2016-10-01"},
		{Amount: 50, CounterPartyName: "Claudia Klar", Usage: "Sparen Samuel", Date: "2016-09-01"},
	}

	testMux.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"amount":50,"counterPartyName":"Claudia Klar","usage":"Sparen Samuel","date":"2016-10-01"},{"amount":50,"counterPartyName":"Claudia Klar","usage":"Sparen Samuel","date":"2016-09-01"}]`)
	})

	act, _, err := testClient.Transactions.Get("DE10000000000000000455")
	ok(t, err)
	equals(t, exp, act)
}
