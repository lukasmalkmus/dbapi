package dbapi

// The AuthenticationService is a simple wrapper around the authentication
// credentials/tokens.
type AuthenticationService struct {
	token string
}

// HasAuth describes if authentication credentials are set.
func (s *AuthenticationService) HasAuth() bool {
	return len(s.token) > 0
}

// Token returns the access token.
func (s *AuthenticationService) Token() string {
	return s.token
}
