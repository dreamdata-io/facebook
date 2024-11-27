package facebook

import (
	"context"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type MeAPI interface {
	AdAccounts(ctx context.Context, params Params) (Result, error)
	Me(ctx context.Context, params Params) (Result, error)
	User(ctx context.Context) (User, error)
}

func (c *Client) AdAccounts(ctx context.Context, params Params) (Result, error) {
	return c.session.Get("/me/adaccounts", params)
}

func (c *Client) Me(ctx context.Context, params Params) (Result, error) {
	res, err := c.session.Get("/me", params)
	if err != nil {
		return Result{}, err
	}

	return res, nil
}

func (c *Client) User(_ context.Context) (User, error) {
	res, err := c.Me(context.Background(), FieldsParams("id", "email"))
	if err != nil {
		return User{}, err
	}

	var user User
	err = res.Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
