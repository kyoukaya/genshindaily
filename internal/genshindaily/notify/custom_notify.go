package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
)

type CustomNotifier struct {
	URL     string            `json:"url" validate:"url"`
	Method  string            `json:"http_method" validate:"oneof=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	Headers map[string]string `json:"http_headers"`
	Body    string            `json:"body"`

	ExpectHTTPCode int `json:"expect_http_code" validate:"min=100,max=599"`
}

func NewCustomNotifier(msg json.RawMessage) (Notifier, error) {
	var c *CustomNotifier
	return c, json.Unmarshal(msg, c)
}

func (c *CustomNotifier) Name() string {
	return "Custom Notifier"
}

func (c *CustomNotifier) Do(ctx context.Context, cli *resty.Client, res *models.Result) error {
	r := cli.R().SetContext(ctx)
	if c.Body != "" {
		r = r.SetBody(c.Body)
	}
	if len(c.Headers) > 0 {
		r = r.SetHeaders(c.Headers)
	}
	resp, err := r.Execute(c.Method, c.URL)
	if err != nil {
		return err
	}
	if code := resp.StatusCode(); code != c.ExpectHTTPCode {
		return fmt.Errorf("http code received %d != %d", code, c.ExpectHTTPCode)
	}
	return nil
}
