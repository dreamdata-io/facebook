package facebook

import (
	"context"
	"fmt"
)

type ConversionsAPI interface {
	Dataset(ctx context.Context, datasetID string, params Params) (Result, error)
	Datasets(ctx context.Context, adAccountID string, params Params) (Result, error)
	UploadEvents(ctx context.Context, datasetID string, params Params) (Result, error)
}

func (c *Client) Dataset(ctx context.Context, datasetID string, params Params) (Result, error) {
	return c.session.Get(fmt.Sprintf("/%s", datasetID), params)
}

func (c *Client) Datasets(ctx context.Context, adAccountID string, params Params) (Result, error) {
	return c.session.Get(fmt.Sprintf("/%s/adspixels", adAccountID), params)
}

func (c *Client) UploadEvents(ctx context.Context, datasetID string, params Params) (Result, error) {
	return c.session.Post(fmt.Sprintf("/%s/events", datasetID), params)
}
