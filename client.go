package facebook

import (
	"context"
	"fmt"
	"github.com/dreamdata-io/facebook/internal"
	"golang.org/x/oauth2"
)

type IClient interface {
	Auth(context.Context, *oauth2.Token, ...AuthOption) IClient
	Session() *internal.Session
	AuthClient
	ConversionsAPI
	MeAPI
	AudiencesAPI
}

type Client struct {
	oauth2Config *oauth2.Config
	version      string
	session      *internal.Session
	app          *internal.App
}

var _ IClient = (*Client)(nil)

func New(cfg Config) *Client {
	app := internal.New(cfg.OAuth2.ClientID, cfg.OAuth2.ClientSecret)
	app.RedirectUri = cfg.OAuth2.RedirectURL

	return &Client{
		version: cfg.Version,
		oauth2Config: &oauth2.Config{
			ClientID:     cfg.OAuth2.ClientID,
			ClientSecret: cfg.OAuth2.ClientSecret,
			RedirectURL:  cfg.OAuth2.RedirectURL,
			Scopes:       cfg.OAuth2.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("https://www.facebook.com/%s/dialog/oauth", cfg.Version),
				TokenURL: fmt.Sprintf("https://graph.facebook.com/%s/oauth/access_token", cfg.Version),
			},
		},
		app: internal.New(cfg.OAuth2.ClientID, cfg.OAuth2.ClientSecret),
	}
}

func (c *Client) Session() *internal.Session {
	return c.session
}
