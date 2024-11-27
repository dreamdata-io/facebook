package facebook

import (
	"context"
	"fmt"
)

type ConversionsAPI interface {
	Dataset(ctx context.Context, datasetID string, params Params) (Result, error)
}

func (c *Client) Dataset(ctx context.Context, datasetID string, params Params) (Result, error) {
	return c.session.Get(fmt.Sprintf("/%s", datasetID), params)
}
