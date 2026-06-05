package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
)

func StartWorker(ctx context.Context, uc *CheckoutUsecase, consumer domain.QueueConsumer, logger *util.Logger) {
	slog.Info("Starting checkout worker")

	go func() {
		backoff := time.Second
		const maxBackoff = 30 * time.Second

		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopping checkout worker")
				return
			default:
			}

			deliveries, err := consumer.Consume()
			if err != nil {
				slog.Error("Failed to start consuming", slog.String("error", err.Error()), slog.Duration("retry_in", backoff))
				select {
				case <-ctx.Done():
					return
				case <-time.After(backoff):
				}
				if backoff < maxBackoff {
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
				continue
			}

			backoff = time.Second // reset on successful connection

			for delivery := range deliveries {
				var job CheckoutJob
				err := json.Unmarshal(delivery.Body, &job)
				if err != nil {
					slog.Error("Failed to unmarshal job", slog.String("error", err.Error()))
					delivery.Nack(false, false)
					continue
				}

				err = uc.ProcessJob(ctx, job)
				if err != nil {
					slog.Error("Failed to process job", slog.String("error", err.Error()))
				}

				delivery.Ack(false)
			}
		}
	}()
}
