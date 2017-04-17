package dbapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUserInfoService_Get(t *testing.T) {
	setup()
	defer teardown()

	exp := &UserInfo{FirstName: "Claudia", LastName: "Klar", Gender: "FEMALE", DateOfBirth: "1977-03-02"}

	testMux.HandleFunc("/v1/userInfo", func(w http.ResponseWriter, r *http.Request) {
		equals(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{"dateOfBirth":"1977-03-02","firstName":"Claudia","gender":"FEMALE","lastName":"Klar"}`)
	})

	act, _, err := testClient.UserInfo.Get()
	ok(t, err)
	equals(t, exp, act)
}
