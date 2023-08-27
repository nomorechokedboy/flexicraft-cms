package config

type GitHubOauth struct {
	GitHubClientID         string `env:"GITHUB_OAUTH_CLIENT_ID"`
	GitHubClientSecret     string `env:"GITHUB_OAUTH_CLIENT_SECRET"`
	GitHubOAuthRedirectUrl string `env:"GITHUB_OAUTH_REDIRECT_URL"`
}
