package notify

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
)

func TestInitNotifier(t *testing.T) {
	res := &models.Result{}
	type mock struct {
		url       string
		method    string
		responder httpmock.Responder
	}
	tests := []struct {
		name        string
		msgFilename string
		mock        *mock
		wantErr     bool
	}{
		{
			name:        "Discord",
			msgFilename: "testdata/discord_example.json",
			mock: &mock{
				url:       "https://discord.com/api/webhooks/123456789012345678/bar",
				method:    http.MethodPost,
				responder: httpmock.NewStringResponder(http.StatusNoContent, ""),
			},
		},
		{
			name:        "healthcheck.io",
			msgFilename: "testdata/health_check_example.json",
			mock: &mock{
				url:       "https://hc-ping.com/foo",
				method:    http.MethodGet,
				responder: httpmock.NewStringResponder(http.StatusOK, ""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				httpmock.ActivateNonDefault(http.DefaultClient)
				defer httpmock.Deactivate()
				httpmock.RegisterResponder(tt.mock.method, tt.mock.url, tt.mock.responder)
			}
			b, err := ioutil.ReadFile(tt.msgFilename)
			assert.NilError(t, err)
			got, err := InitNotifier(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitNotifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NilError(t, got.Do(context.Background(), resty.NewWithClient(http.DefaultClient), res))
		})
	}
}
