package vk

import (
	"github.com/pkg/errors"
	"time"
)

type Poller interface {
	Poll(b *Bot) error
}

type LongPoller struct {
	Wait            int
	RequestInterval time.Duration
}

func (l *LongPoller) Poll(b *Bot) error {
	server, err := b.GetLongPollServer(l.Wait)
	if err != nil {
		return errors.Wrap(err, "Ошибка запуска Long Poller`a")
	}

	b.Logger.Info("Начинаем Long Polling...")

	for {
		updates, err := GetUpdates(b, server)
		if err != nil {
			b.Logger.Warnf("[LongPoller] Ошибка получния апдейтов: %+v", err)
			time.Sleep(l.RequestInterval)
			continue
		}

		b.ProcessUpdates(updates)
		time.Sleep(l.RequestInterval)
	}

}

func NewDefaultLongPoller() *LongPoller {
	return &LongPoller{
		Wait:            DEFAULT_LONG_POLLING_WAIT,
		RequestInterval: time.Millisecond * DEFAULT_REQUEST_INTERVAL,
	}
}
