package api

import (
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/oauth2"
)

const (
	authServerURL = "http://localhost:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "test_client_1",
		ClientSecret: "test_secret_1",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8088/index",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
