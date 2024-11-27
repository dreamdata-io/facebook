package facebook

import (
	"context"
	"errors"
	"github.com/dreamdata-io/facebook/internal"
	"golang.org/x/oauth2"
)

type AuthClient interface {
	ClientID() string
	ClientSecret() string
	OAuth2Config() *oauth2.Config
	AuthCodeURL(ctx context.Context, state string, options ...oauth2.AuthCodeOption) string
	ExchangeOAuth2Code(ctx context.Context, oauth2Code string) (*oauth2.Token, error)
	AccessToken(ctx context.Context, refreshToken string) (*oauth2.Token, error)
	// Revoke revokes the token, this method requires authorization
	Revoke(ctx context.Context, refreshToken string) error
}

type AuthOption func(*oauth2.Config)

func WithScopes(scopes ...string) AuthOption {
	return func(cfg *oauth2.Config) {
		cfg.Scopes = scopes
	}
}

func (c *Client) Auth(ctx context.Context, token *oauth2.Token, opts ...AuthOption) IClient {
	cfg := c.oauth2Config
	for _, option := range opts {
		option(cfg)
	}

	session := c.app.Session("")
	session.Version = c.version
	session.HttpClient = cfg.Client(ctx, token)

	return &Client{
		app:          c.app,
		oauth2Config: cfg,
		session:      session,
		version:      c.version,
	}
}

func (c *Client) ClientID() string {
	return c.oauth2Config.ClientID
}

func (c *Client) ClientSecret() string {
	return c.oauth2Config.ClientSecret
}

func (c *Client) OAuth2Config() *oauth2.Config {
	return c.oauth2Config
}

func (c *Client) AuthCodeURL(_ context.Context, state string, options ...oauth2.AuthCodeOption) string {
	return c.oauth2Config.AuthCodeURL(state, options...)
}

func (c *Client) ExchangeOAuth2Code(ctx context.Context, oauth2Code string) (*oauth2.Token, error) {
	return c.oauth2Config.Exchange(ctx, oauth2Code)
}

func (c *Client) AccessToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	t, err := c.oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Client) Revoke(_ context.Context, refreshToken string) error {
	params := map[string]string{
		"client_id":     c.oauth2Config.ClientID,
		"client_secret": c.oauth2Config.ClientSecret,
		"revoke_token":  refreshToken,
	}

	res, err := c.session.Get("/oauth/revoke", internal.MakeParams(params))
	if err != nil {
		return err
	}

	var revocation struct {
		Success bool `json:"success"`
	}

	if err = res.Decode(&revocation); err != nil {
		return err
	}

	if !revocation.Success {
		return errors.New("failed to revoke token")
	}

	return nil
}
