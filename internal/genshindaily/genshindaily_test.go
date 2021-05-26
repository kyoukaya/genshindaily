// package genshindaily contains a Pub/Sub Cloud Function.
package genshindaily

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
	"github.com/kyoukaya/genshindaily/internal/genshindaily/notify"
)

type mock struct {
	url       string
	method    string
	responder httpmock.Responder
}

var res = &models.Result{
	Today: "2021-05-23",
	Award: models.Award{
		Icon: "https://uploadstatic-sea.mihoyo.com/event/2021/02/25/01ba12730bd86c8858c1e2d86c7d150d_5665148762126820826.png",
		Cnt:  5,
		Name: "Adventurer's Experience",
	},
	DaysCheckedIn: 22,
	Status:        models.CheckInStatusDupe,
	UID:           "390XXXXX",
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name      string
		cookie    string
		httpMocks []mock
		expect    *models.Result
		notifiers []string
		wantErr   bool
	}{
		// TODO: Add more test cases
		{
			name:   "Duplicate sign in",
			cookie: "account_id=390XXXXX; cookie_token=XXXXXX",
			httpMocks: []mock{
				{
					method:    http.MethodGet,
					responder: httpmock.NewBytesResponder(http.StatusOK, loadFile(t, "testdata/sign_info_nominal.json")),
					url:       "https://hk4e-api-os.mihoyo.com/event/sol/info?lang=en-us&act_id=e202102251931481",
				},
				{
					method:    http.MethodGet,
					responder: httpmock.NewBytesResponder(http.StatusOK, loadFile(t, "testdata/rewards_info_nominal.json")),
					url:       "https://hk4e-api-os.mihoyo.com/event/sol/home?lang=en-us&act_id=e202102251931481",
				},
				// {
				// 	method:          http.MethodPost,
				// 	payloadFilename: "testdata/sign_in_dupe.json",
				// 	url:             "https://hk4e-api-os.mihoyo.com/event/sol/sign?lang=en-us",
				// },
			},
			expect: res,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.ActivateNonDefault(http.DefaultClient)
			for _, h := range tt.httpMocks {
				httpmock.RegisterResponder(h.method, h.url, h.responder)
			}
			defer httpmock.Deactivate()
			h := &DailyHandler{
				r: resty.NewWithClient(http.DefaultClient),
			}
			got, err := h.SignIn(parseCookies(tt.cookie))
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(got, tt.expect); diff != "" {
				t.Errorf("SignIn() output differs from expected: %s", diff)
			}
		})
	}
}

func TestDailyHandler_Notify(t *testing.T) {
	tests := []struct {
		name      string
		notifiers []notify.Notifier
		mocks     []mock
		res       *models.Result
		wantErr   bool
	}{
		{
			name: "Discord notification failure",
			notifiers: []notify.Notifier{&notify.DiscordNotifier{
				URL: "https://discord.com/api/webhooks/123456789012345678/bar",
			}},
			mocks: []mock{{
				url:       "https://discord.com/api/webhooks/123456789012345678/bar",
				method:    http.MethodPost,
				responder: httpmock.NewStringResponder(http.StatusBadRequest, ""),
			}},
			res:     res,
			wantErr: true,
		},
		{
			name: "Discord multi notification failure",
			notifiers: []notify.Notifier{&notify.DiscordNotifier{
				URL: "https://discord.com/api/webhooks/123456789012345678/bar",
			}, &notify.DiscordNotifier{
				URL: "https://discord.com/api/webhooks/123456789012345678/foo"}},
			mocks: []mock{{
				url:       "https://discord.com/api/webhooks/123456789012345678/bar",
				method:    http.MethodPost,
				responder: httpmock.NewStringResponder(http.StatusBadRequest, ""),
			}, {
				url:       "https://discord.com/api/webhooks/123456789012345678/foo",
				method:    http.MethodPost,
				responder: httpmock.NewStringResponder(http.StatusBadRequest, ""),
			}},
			res:     res,
			wantErr: true,
		},
		{
			name:      "No notifications",
			notifiers: nil,
			res:       res,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.Deactivate()
			h := &DailyHandler{notifiers: tt.notifiers}
			err := h.Notify(context.Background(), tt.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("DailyHandler.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func loadFile(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	assert.NilError(t, err)
	return b
}
