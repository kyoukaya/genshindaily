package genshindaily

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-multierror"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
	"github.com/kyoukaya/genshindaily/internal/genshindaily/notify"
)

// Compile-time variables

var Version = "dev"
var Date = "1 January 1970 12:00:00"
var Commit = "ffffffff"

type Config struct {
	Cookies   string            `json:"cookies" validate:"required"`
	Notifiers []json.RawMessage `json:"notifiers"`
}

func HandleMessage(ctx context.Context, c *Config) (string, error) {
	h := &DailyHandler{
		r:      resty.NewWithClient(http.DefaultClient),
		Config: c,
	}
	res, err := h.HandleMessage(ctx)
	if err != nil {
		return "", err
	}
	return res.String(), err
}

type DailyHandler struct {
	r *resty.Client
	*Config
	notifiers []notify.Notifier
}

func (h *DailyHandler) HandleMessage(ctx context.Context) (*models.Result, error) {
	t0 := time.Now()
	fmt.Printf("genshindaily %s (%s)\nBuilt: %s\n\n", Version, Commit, Date)

	if err := h.InitNotifiers(ctx); err != nil {
		return nil, err
	}

	res, err := h.SignIn(parseCookies(h.Cookies))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Signing in OK, status: %s. Took %.2fs\n", res.Status, time.Since(t0).Seconds())
	res.Start = t0
	res.VersionText = Version

	if errs := h.Notify(ctx, res); errs != nil {
		fmt.Printf("notify failed: %v", errs)
		return nil, errs
	}

	fmt.Printf("Total time: %.2fs\n", time.Since(t0).Seconds())
	return res, nil
}

func (h *DailyHandler) InitNotifiers(ctx context.Context) error {
	fmt.Printf("initializing %d notifiers\n", len(h.Notifiers))
	for _, msg := range h.Notifiers {
		n, err := notify.InitNotifier(msg)
		if err != nil {
			return err
		}
		h.notifiers = append(h.notifiers, n)
	}
	return nil
}

func (h *DailyHandler) Notify(ctx context.Context, res *models.Result) error {
	t1 := time.Now()
	errch := make(chan error, len(h.notifiers))
	wg := &sync.WaitGroup{}
	for i := range h.notifiers {
		n := h.notifiers[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			errch <- n.Do(ctx, resty.NewWithClient(http.DefaultClient), res)
		}()
	}
	wg.Wait()
	close(errch)

	var errs error
	for err := range errch {
		if err == nil {
			continue
		}
		errs = multierror.Append(errs, err)
	}
	fmt.Printf("Notifications took %.2fs\n", time.Since(t1).Seconds())
	return errs
}
