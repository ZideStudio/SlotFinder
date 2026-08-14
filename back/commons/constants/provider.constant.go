package constants

type Provider string

const (
	PROVIDER_GOOGLE  Provider = "google"
	PROVIDER_DISCORD Provider = "discord"
	PROVIDER_GITHUB  Provider = "github"
)

var PROVIDERS = []Provider{
	PROVIDER_GOOGLE,
	PROVIDER_DISCORD,
	PROVIDER_GITHUB,
}
