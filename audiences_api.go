package facebook

import (
	"context"
	"fmt"
	"github.com/dreamdata-io/facebook/internal"
	"strings"
)

type AudiencesAPI interface {
	Audience(ctx context.Context, audienceID string, params Params) (Result, error)
	CustomAudiences(ctx context.Context, adAccountID string, params Params) (Result, error)
	CreateAudience(ctx context.Context, adAccountID string, params Params) (Result, error)
}

// Audience calls the Facebook Graph API with GET at /{audience_id} to get an audience.
func (c *Client) Audience(_ context.Context, audienceID string, params Params) (Result, error) {
	return c.session.Get(fmt.Sprintf("/%s", audienceID), params)
}

// CustomAudiences calls the Facebook Graph API with GET at /act_{ad_account_id}/customaudiences to get all audiences.
func (c *Client) CustomAudiences(_ context.Context, adAccountID string, params Params) (Result, error) {
	return c.session.Get(fmt.Sprintf("/act_%s/customaudiences", adAccountID), internal.MakeParams(params))
}

type AudienceSubtype = string

type FileSource = string

const (
	CustomAudience AudienceSubtype = "CUSTOM"

	UserProvidedOnlyFileSource     FileSource = "USER_PROVIDED_ONLY"
	PartnerProvidedOnlyFileSource  FileSource = "PARTNER_PROVIDED_ONLY"
	BothUsersAndPartnersFileSource FileSource = "BOTH_USERS_AND_PARTNERS_PROVIDED"
)

// CreateAudience calls the Facebook API with POST at /act_{ad_account_id}/customaudiences to create a new audience.
func (c *Client) CreateAudience(ctx context.Context, adAccountID string, params Params) (Result, error) {
	return c.session.Post(fmt.Sprintf("/act_%s/customaudiences", adAccountID), params)
}

type AddUserSession struct {
	SessionID         string `json:"session_id"`
	BatchSeq          int    `json:"batch_seq"`
	LastBatchFlag     bool   `json:"last_batch_flag"`
	EstimatedNumTotal int    `json:"estimated_num_total"`
}

type AddUserPayload struct {
	Schema []string         `json:"schema"`
	Data   []map[string]any `json:"data"`
}

func (p AddUserPayload) Format() map[string]interface{} {
	var payloadData [][]any
	for _, entry := range p.Data {
		var payloadEntry []any
		for _, key := range p.Schema {
			payloadEntry = append(payloadEntry, entry[strings.ToLower(key)])
		}
		payloadData = append(payloadData, payloadEntry)
	}
	return map[string]interface{}{
		"schema": fmt.Sprintf("%v", p.Schema),
		"data":   fmt.Sprintf("%v", payloadData),
	}
}

type AddUserOutput struct {
	AudienceId          string            `json:"audience_id"`
	SessionId           string            `json:"session_id"`
	NumReceived         int               `json:"num_received"`
	NumInvalidEntries   int               `json:"num_invalid_entries"`
	InvalidEntrySamples map[string]string `json:"invalid_entry_samples"`
}

// AddUsers calls the Facebook API at /{audience_id}/users to add users to an audience.
func (c *Client) AddUsers(ctx context.Context, audienceID string, payload AddUserPayload, params Params) (Result, error) {
	if params == nil {
		params = make(Params)
	}
	params["payload"] = payload.Format()
	return c.session.Post(fmt.Sprintf("/%s/users", audienceID), params)
}

func (c *Client) ReplaceUsers(ctx context.Context, audienceID string, payload AddUserPayload, params Params) (Result, error) {
	res, err := c.Audience(ctx, audienceID, FieldsParams("operation_status"))
	if err != nil {
		return nil, err
	}

	var opStatus int
	err = res.DecodeField("operation_status", &opStatus)
	if err != nil {
		return nil, err
	}
	if opStatus != 200 {
		return nil, fmt.Errorf("audience %s is not ready for replacement", audienceID)
	}

	if params == nil {
		params = make(Params)
	}
	params["payload"] = payload.Format()
	return c.session.Post(fmt.Sprintf("/%s/usersreplace", audienceID), params)
}

// Sessions calls the Facebook API with GET at /{audience_id}/sessions to get information on sessions.
func (c *Client) Sessions(ctx context.Context, audienceID string, sessionID string) (Result, error) {
	return c.session.Get(fmt.Sprintf("/%s/sessions", audienceID), Params{"session_id": sessionID})
}
