package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/renantatsuo/app-review/server/internal/apps"
	"github.com/renantatsuo/app-review/server/internal/config"
	"github.com/renantatsuo/app-review/server/internal/queue"
)

type Scheduler struct {
	l          *slog.Logger
	appsClient *apps.AppsClient
	queue      queue.Queue
	config     config.Config
}

func New(l *slog.Logger, appsClient *apps.AppsClient, queue queue.Queue, config config.Config) *Scheduler {
	return &Scheduler{l: l, appsClient: appsClient, queue: queue, config: config}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.l.Info("starting scheduler")

	ticker := time.NewTicker(s.config.PollingInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				apps, err := s.appsClient.GetAllApps()
				if err != nil {
					s.l.Error("error getting all apps", "error", err)
					continue
				}

				for _, app := range apps {
					s.l.Info("scheduling app", "app", app)
					if err := s.queue.Enqueue([]byte(app.ID)); err != nil {
						s.l.Error("error enqueuing app", "error", err)
						continue
					}
				}
			}
		}
	}()
}
