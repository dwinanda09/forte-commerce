package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/dwinanda09/forte-commerce/util"
)

func StartExpiryWorker(ctx context.Context, uc *CheckoutUsecase, logger *util.Logger) {
	slog.Info("Starting expiry worker")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("Stopping expiry worker")
				return
			case <-ticker.C:
				err := uc.ReleaseExpired(ctx)
				if err != nil {
					slog.Error("Failed to release expired sessions", slog.String("error", err.Error()))
				}
			}
		}
	}()
}
