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

// Discord
const (
	PROVIDER_DISCORD_URL          = "https://discord.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=identify+email&state=%s"
	PROVIDER_DISCORD_TOKEN_URL    = "https://discord.com/api/oauth2/token"
	PROVIDER_DISCORD_USERINFO_URL = "https://discord.com/api/users/@me"
)

// Google
const (
	PROVIDER_GOOGLE_URL          = "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email%%20profile&state=%s"
	PROVIDER_GOOGLE_TOKEN_URL    = "https://oauth2.googleapis.com/token"
	PROVIDER_GOOGLE_USERINFO_URL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// Github
const (
	PROVIDER_GITHUB_URL           = "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s"
	PROVIDER_GITHUB_TOKEN_URL     = "https://github.com/login/oauth/access_token"
	PROVIDER_GITHUB_USERINFO_URL  = "https://api.github.com/user"
	PROVIDER_GITHUB_USEREMAIL_URL = "https://api.github.com/user/emails"
)
