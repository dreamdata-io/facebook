package facebook

type (
	Config struct {
		Version string       `envconfig:"VERSION" default:"v21.0"`
		OAuth2  OAuth2Config `envconfig:"OAUTH2"`
	}

	OAuth2Config struct {
		ClientID     string   `envconfig:"CLIENT_ID" required:"true"`
		ClientSecret string   `envconfig:"CLIENT_SECRET" required:"true"`
		Scopes       []string `envconfig:"SCOPES"`
		RedirectURL  string   `envconfig:"REDIRECT_URL"`
	}
)
