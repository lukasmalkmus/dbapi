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
		{OriginIBAN: "DE10000000000000000454", Amount: -35.56, CounterPartyName: "Netto", Usage: "POS MIT PIN. Einkauf", BookingDate: "2016-10-27"},
		{OriginIBAN: "DE10000000000000000454", Amount: -52.22, CounterPartyName: "Lidl", Usage: "POS MIT PIN. Einkauf", BookingDate: "2016-10-24"},
		{OriginIBAN: "DE10000000000000000454", Amount: -1500, CounterPartyName: "Schwäbisch Hall", Usage: "Ref. 58974-8765889", BookingDate: "2016-10-21"},
		{OriginIBAN: "DE10000000000000000454", Amount: -38.98, CounterPartyName: "Toys R Us", Usage: "Rechnung", BookingDate: "2016-10-17"},
		{OriginIBAN: "DE10000000000000000454", Amount: -25.95, CounterPartyName: "Alnatura Frankfurt", Usage: "POS MIT PIN. Einkauf", BookingDate: "2016-10-17"},
		{OriginIBAN: "DE10000000000000000454", Amount: -96.16, CounterPartyName: "JET", Usage: "POS MIT PIN. Die Tanke Ihrer Wahl", BookingDate: "2016-10-12"},
	}

	testMux.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"originIBAN":"DE10000000000000000454","amount":-35.56,"counterPartyName":"Netto","usage":"POS MIT PIN. Einkauf","bookingDate":"2016-10-27"},{"originIBAN":"DE10000000000000000454","amount":-52.22,"counterPartyName":"Lidl","usage":"POS MIT PIN. Einkauf","bookingDate":"2016-10-24"},{"originIBAN":"DE10000000000000000454","amount":-1500,"counterPartyName":"Schwäbisch Hall","usage":"Ref. 58974-8765889","bookingDate":"2016-10-21"},{"originIBAN":"DE10000000000000000454","amount":-38.98,"counterPartyName":"Toys R Us","usage":"Rechnung","bookingDate":"2016-10-17"},{"originIBAN":"DE10000000000000000454","amount":-25.95,"counterPartyName":"Alnatura Frankfurt","usage":"POS MIT PIN. Einkauf","bookingDate":"2016-10-17"},{"originIBAN":"DE10000000000000000454","amount":-96.16,"counterPartyName":"JET","usage":"POS MIT PIN. Die Tanke Ihrer Wahl","bookingDate":"2016-10-12"}]`)
	})

	act, _, err := testClient.Transactions.GetAll()
	ok(t, err)
	equals(t, exp, act)
}

func TestTransactionsService_Get(t *testing.T) {
	setup()
	defer teardown()

	exp := &Transactions{
		{OriginIBAN: "DE10000000000000000455", Amount: 50, CounterPartyName: "Claudia Klar", Usage: "Sparen Samuel", BookingDate: "2016-10-01"},
		{OriginIBAN: "DE10000000000000000455", Amount: 50, CounterPartyName: "Claudia Klar", Usage: "Sparen Samuel", BookingDate: "2016-09-01"},
	}

	testMux.HandleFunc("/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"originIBAN":"DE10000000000000000455","amount":50,"counterPartyName":"Claudia Klar","usage":"Sparen Samuel","bookingDate":"2016-10-01"},{"originIBAN":"DE10000000000000000455","amount":50,"counterPartyName":"Claudia Klar","usage":"Sparen Samuel","bookingDate":"2016-09-01"}]`)
	})

	act, _, err := testClient.Transactions.Get("DE10000000000000000455")
	ok(t, err)
	equals(t, exp, act)
}
