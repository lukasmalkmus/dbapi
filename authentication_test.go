package dbapi

import "testing"

func TestHasAuth(t *testing.T) {
	mockData := []struct {
		token          string
		ExpectedStatus bool
		ExpectedError  error
	}{
		{"", false, nil},
		{"123", true, nil},
		{testAccessToken, true, nil},
	}

	for _, mock := range mockData {
		c, err := NewClient(
			SetToken(mock.token),
		)
		if c != nil {
			equals(t, mock.ExpectedStatus, c.Authentication.HasAuth())
		}
		equals(t, err, mock.ExpectedError)
	}
}

func TestToken(t *testing.T) {
	mockData := []struct {
		token         string
		ExpectedToken string
		ExpectedError error
	}{
		{"", "", nil},
		{"123", "123", nil},
		{testAccessToken, testAccessToken, nil},
	}

	for _, mock := range mockData {
		c, err := NewClient(
			SetToken(mock.token),
		)
		if c != nil {
			equals(t, mock.ExpectedToken, c.Authentication.Token())
		}
		equals(t, err, mock.ExpectedError)
	}
}
