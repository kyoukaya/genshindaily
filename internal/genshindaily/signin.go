package genshindaily

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
)

const (
	DefaultLang = "en-us"
	ActID       = "e202102251931481"
	UserAgent   = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"
	RefererURL  = "https://webstatic-sea.mihoyo.com/ys/event/signin-sea/index.html?act_id=%s"
	InfoURL     = "https://hk4e-api-os.mihoyo.com/event/sol/info?lang=%s&act_id=%s"
	RewardURL   = "https://hk4e-api-os.mihoyo.com/event/sol/home?lang=%s&act_id=%s"
	SignURL     = "https://hk4e-api-os.mihoyo.com/event/sol/sign?lang=%s"
)

func (h *DailyHandler) SignIn(cookies []*http.Cookie) (*models.Result, error) {
	uid := getCookieByKey(cookies, "account_id")
	token := getCookieByKey(cookies, "cookie_token")
	if uid == nil || token == nil {
		return nil, errors.New("account_id or cookie_token is not provided")
	}

	r := h.r.
		SetHeaders(map[string]string{
			"User-Agent":      UserAgent,
			"Referer":         fmt.Sprintf(RefererURL, ActID),
			"Accept-Encoding": "gzip, deflate, br",
		}).
		SetCookies(cookies)

	signInf, rewards, err := h.GetInfoAndRewards()
	if err != nil {
		return nil, err
	}

	res := &models.Result{
		UID:           uid.Value,
		Today:         signInf.Today,
		DaysCheckedIn: signInf.TotalSignDay,
	}

	// Already signed in today.
	if signInf.IsSign {
		res.Award = rewards.Awards[signInf.TotalSignDay-1]
		res.Status = models.CheckInStatusDupe
		return res, nil
	}

	res.Award = rewards.Awards[signInf.TotalSignDay]
	// Honestly don't what this is, but will mirror the logic.
	if signInf.FirstBind {
		res.Status = models.CheckInStatusFirstBind
		return res, nil
	}

	err = mhyRequestWrapper(
		r.R().SetBody(&models.SignInRequest{ActID: ActID}),
		http.MethodPost,
		fmt.Sprintf(SignURL, DefaultLang),
		nil)
	if err == ErrAlreadyCheckedIn {
		res.Status = models.CheckInStatusDupe
		return res, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to sign in: %v", err)
	}

	res.DaysCheckedIn++
	res.Status = models.CheckInStatusOK
	return res, nil
}

func (h *DailyHandler) GetInfoAndRewards() (info *models.InfoResponse, rewards *models.RewardsResponse, err error) {
	g, ctx := errgroup.WithContext(context.TODO())
	g.Go(func() error { return GetInfo(ctx, h.r, &info) })
	g.Go(func() error { return GetRewards(ctx, h.r, &rewards) })
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return info, rewards, nil
}

// Get Info Response
func GetInfo(ctx context.Context, r *resty.Client, info **models.InfoResponse) error {
	return mhyRequestWrapper(r.R(), http.MethodGet, fmt.Sprintf(InfoURL, DefaultLang, ActID), info)
}

func GetRewards(ctx context.Context, r *resty.Client, rewards **models.RewardsResponse) error {
	return mhyRequestWrapper(r.R(), http.MethodGet, fmt.Sprintf(RewardURL, DefaultLang, ActID), rewards)
}

var ErrAlreadyCheckedIn = errors.New("already checked in")

func mhyRequestWrapper(r *resty.Request, method, url string, resp interface{}) error {
	outerResp := &models.OuterResponse{
		Data: resp,
	}
	httpResp, err := r.Execute(method, url)
	if err != nil {
		return err
	}
	if code := httpResp.StatusCode(); code != http.StatusOK {
		return fmt.Errorf("got non 200 HTTP code: %d", code)
	}
	if err := json.Unmarshal(httpResp.Body(), outerResp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	if code := outerResp.Retcode; code != 0 {
		if code == -5003 {
			return ErrAlreadyCheckedIn
		}
		return fmt.Errorf("got non 0 retcode: %d", outerResp.Retcode)
	}
	return nil
}
