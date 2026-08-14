package provider

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
)

// Real OAuth endpoints called by discord.provider.go / github.provider.go /
// google.provider.go. Tests redirect these exact URLs to local httptest
// servers instead of ProviderService exposing per-URL fields, so the
// production code only gains a single injectable *resty.Client.
const (
	realDiscordTokenURL    = "https://discord.com/api/oauth2/token"
	realDiscordUserInfoURL = "https://discord.com/api/users/@me"

	realGoogleTokenURL    = "https://oauth2.googleapis.com/token"
	realGoogleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	realGithubTokenURL     = "https://github.com/login/oauth/access_token"
	realGithubUserInfoURL  = "https://api.github.com/user"
	realGithubUserEmailURL = "https://api.github.com/user/emails"
)

// redirectTransport rewrites requests whose exact URL matches a key in
// routes to the corresponding local test-server URL, then delegates to next.
type redirectTransport struct {
	routes map[string]string
	next   http.RoundTripper
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	target, ok := rt.routes[req.URL.String()]
	if !ok {
		return rt.next.RoundTrip(req)
	}

	newURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	req = req.Clone(req.Context())
	req.URL = newURL
	req.Host = ""

	return rt.next.RoundTrip(req)
}

// oauthRedirectClient returns a *resty.Client whose requests matching routes
// (real provider URL -> local httptest server URL) are rerouted there;
// anything else is delegated to http.DefaultTransport. resty.New() builds
// its own concrete *http.Transport rather than falling back to
// http.DefaultTransport (see createClient in go-resty/resty/v2@v2.17.2's
// client.go), so overriding the global http.DefaultTransport has no effect
// here — the redirect must be wired into the client actually assigned to
// ProviderService.httpClient.
func oauthRedirectClient(t *testing.T, routes map[string]string) *resty.Client {
	t.Helper()
	return resty.NewWithClient(&http.Client{
		Transport: &redirectTransport{routes: routes, next: http.DefaultTransport},
	})
}
