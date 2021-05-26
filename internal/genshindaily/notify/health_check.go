package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
)

type HealthCheck struct {
	URL string `json:"url" validate:"url"`
}

func NewHealthChecker(msg json.RawMessage) (Notifier, error) {
	h := &HealthCheck{}
	return h, json.Unmarshal(msg, h)
}

func (h *HealthCheck) Name() string { return "Health checker" }

func (h *HealthCheck) Do(ctx context.Context, cli *resty.Client, res *models.Result) error {
	c := &CustomNotifier{
		URL:            h.URL,
		Method:         http.MethodGet,
		ExpectHTTPCode: http.StatusOK,
	}
	if err := c.Do(ctx, cli, res); err != nil {
		return fmt.Errorf("failed to do health check: %w", err)
	}
	return nil
}
