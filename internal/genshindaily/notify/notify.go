package notify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"

	"github.com/kyoukaya/genshindaily/internal/genshindaily/models"
)

type Notifier interface {
	Do(ctx context.Context, cli *resty.Client, res *models.Result) error
}

type baseNotifier struct {
	Kind string `json:"kind"`
}

type NotifierFactory func(msg json.RawMessage) (Notifier, error)

var notifyFactoryMap = map[string]NotifierFactory{
	"discord":      NewDiscordNotifier,
	"health_check": NewHealthChecker,
}

func InitNotifiers(msgs [][]byte) ([]Notifier, error) {
	var notifiers []Notifier
	for i, msg := range msgs {
		n, err := InitNotifier(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to init notifier in index %d: %w", i, err)
		}
		notifiers = append(notifiers, n)
	}
	return notifiers, nil
}

var validate = validator.New()

func InitNotifier(msg []byte) (Notifier, error) {
	base := baseNotifier{}
	err := json.Unmarshal(msg, &base)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	f, ok := notifyFactoryMap[base.Kind]
	if !ok {
		return nil, fmt.Errorf("no notifier of kind %q", base.Kind)
	}
	n, err := f(msg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", base.Kind, err)
	}
	return n, validate.Struct(n)
}
