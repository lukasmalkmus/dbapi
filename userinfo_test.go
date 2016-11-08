package dbapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserInfoService_Get(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/userInfo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, GET)
		fmt.Fprint(w, `{"dateOfBirth":"1977-03-02","firstName":"Claudia","gender":"FEMALE","lastName":"Klar"}`)
	})

	act, _, err := testClient.UserInfo.Get()
	ok(t, err)

	exp := &UserInfo{DateOfBirth: "1977-03-02", FirstName: "Claudia", Gender: "FEMALE", LastName: "Klar"}

	equals(t, exp, act)
}
