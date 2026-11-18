package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"lab13/order-service/internal/common/exceptions"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	client               Doer
	users, notifications string
}

func New(c Doer, users, notifications string) *Client {
	return &Client{client: c, users: strings.TrimRight(users, "/"), notifications: strings.TrimRight(notifications, "/")}
}

func (c *Client) Exists(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.users+"/api/v1/users/"+id, nil)
	if err != nil {
		return exceptions.ErrUnavailable
	}
	res, err := c.client.Do(req)
	if err != nil {
		return exceptions.ErrUnavailable
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return exceptions.ErrInvalidUser
	}
	if res.StatusCode != 200 {
		return exceptions.ErrUnavailable
	}
	var body struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&body); err != nil || !strings.EqualFold(body.Data.ID, id) {
		return exceptions.ErrUnavailable
	}
	return nil
}

func (c *Client) Send(ctx context.Context, userID, message string) error {
	body, err := json.Marshal(struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}{userID, message})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.notifications+"/api/v1/notifications", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return fmt.Errorf("notification service returned %d", res.StatusCode)
	}
	return nil
}
