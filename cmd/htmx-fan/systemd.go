package main

import (
	"context"

	notify "github.com/iguanesolutions/go-systemd/v5/notify"
	watchdog "github.com/iguanesolutions/go-systemd/v5/notify/watchdog"
	"github.com/rs/zerolog"
)

func systemdWatchdog(log *zerolog.Logger, ctx context.Context) {
	go func() {
		w, err := watchdog.New()
		if err != nil {
			log.Info().Err(err).Msg("Failure to create systemd watchdog")
			return
		}

		for {
			select {
			case <-w.NewTicker().C:
				log.Info().Msg("Watchdog ping")
				if err := w.SendHeartbeat(); err != nil {
					log.Info().Err(err).Msg("Failure to send systemd watchdog heartbeat")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func systemdStopping(log *zerolog.Logger, ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("Notifying systemd that service is stopping")
		if err := notify.Stopping(); err != nil {
			log.Info().Err(err).Msg("Failure to notify systemd that service is stopping")
		}
	}()
}
