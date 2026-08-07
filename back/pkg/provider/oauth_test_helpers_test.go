package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newOAuthFailureServer returns a server that fails every request with a 401,
// used to exercise the "OAuth: failed to get token" error branches.
func newOAuthFailureServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)
	return server
}

// closedServerURL starts and immediately closes a server, returning a URL
// guaranteed to refuse connections — used to exercise the network-error
// branches (as opposed to the non-2xx-status branches above).
func closedServerURL(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()
	return url
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// newOAuthSuccessServer mimics Google's token + userinfo endpoints at
// {url}/token and {url}/userinfo.
func newOAuthSuccessServer(t *testing.T, sub, username, email string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"sub": sub, "name": username, "email": email, "picture": ""})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func setGoogleURLs(tokenURL, userInfoURL string) func() {
	originalToken, originalUserInfo := googleTokenURL, googleUserInfoURL
	googleTokenURL, googleUserInfoURL = tokenURL, userInfoURL
	return func() { googleTokenURL, googleUserInfoURL = originalToken, originalUserInfo }
}

// newDiscordOAuthServer mimics Discord's token + userinfo endpoints.
func newDiscordOAuthServer(t *testing.T, id, username, email string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"id": id, "username": username, "email": email, "avatar": "avatarhash"})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func setDiscordURLs(tokenURL, userInfoURL string) func() {
	originalToken, originalUserInfo := discordTokenURL, discordUserInfoURL
	discordTokenURL, discordUserInfoURL = tokenURL, userInfoURL
	return func() { discordTokenURL, discordUserInfoURL = originalToken, originalUserInfo }
}

// newGithubOAuthServer mimics Github's token + user + emails endpoints.
func newGithubOAuthServer(t *testing.T, id int32, login, primaryEmail string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"access_token": "test-access-token"})
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": id, "login": login, "avatar_url": ""})
	})
	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []map[string]any{
			{"email": "secondary@example.com", "primary": false},
			{"email": primaryEmail, "primary": true},
		})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func setGithubURLs(tokenURL, userInfoURL, emailURL string) func() {
	originalToken, originalUserInfo, originalEmail := githubTokenURL, githubUserInfoURL, githubUserEmailURL
	githubTokenURL, githubUserInfoURL, githubUserEmailURL = tokenURL, userInfoURL, emailURL
	return func() {
		githubTokenURL, githubUserInfoURL, githubUserEmailURL = originalToken, originalUserInfo, originalEmail
	}
}
