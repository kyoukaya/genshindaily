package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"

	"github.com/go-resty/resty/v2"
)

type DiscordNotifier struct {
	URL string `json:"url" validate:"url"`
}

func NewDiscordNotifier(msg json.RawMessage) (Notifier, error) {
	d := &DiscordNotifier{}
	return d, json.Unmarshal(msg, d)
}

func (d *DiscordNotifier) Name() string { return "Discord Webhook Notifier" }

func (d *DiscordNotifier) Do(ctx context.Context, cli *resty.Client, res *models.Result) error {
	b, err := json.Marshal(&DiscordWebhook{Embeds: []*Embed{AsDiscordEmbed(res)}})
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}
	c := &CustomNotifier{
		Method:         http.MethodPost,
		ExpectHTTPCode: http.StatusNoContent,
		Headers:        map[string]string{"Content-Type": "application/json"},
		Body:           string(b),
		URL:            d.URL,
	}
	if err := c.Do(ctx, cli, res); err != nil {
		return fmt.Errorf("failed to notify with discord webhook: %w", err)
	}
	return nil
}

func AsDiscordEmbed(r *models.Result) *Embed {
	return &Embed{
		Title: "genshindaily",
		Color: 0x3b2f8,
		Thumbnail: Image{
			URL: r.Award.Icon,
		},
		Fields: []Field{
			{
				Name:   "UID",
				Inline: true,
				Value:  r.UID,
			},
			{
				Name:   "Status",
				Inline: true,
				Value:  r.Status.String(),
			},
			{
				Name:   "Days Checked In",
				Inline: true,
				Value:  fmt.Sprint(r.DaysCheckedIn),
			},
			{
				Name:   "Today's Reward",
				Inline: true,
				Value:  fmt.Sprintf("%s x %d", r.Award.Name, r.Award.Cnt),
			},
		},
		Footer: Footer{
			Text: fmt.Sprintf("%.2fs - %s",
				time.Since(r.Start).Seconds(),
				r.VersionText,
			),
		},
	}
}

type DiscordWebhook struct {
	Username  string   `json:"username"`
	AvatarURL string   `json:"avatar_url"`
	Content   string   `json:"content"`
	Embeds    []*Embed `json:"embeds"`
}

type Embed struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int64   `json:"color"`
	Fields      []Field `json:"fields"`
	Thumbnail   Image   `json:"thumbnail"`
	Image       Image   `json:"image"`
	Footer      Footer  `json:"footer"`
}

type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type Image struct {
	URL string `json:"url"`
}
