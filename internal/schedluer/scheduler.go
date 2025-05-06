package scheduler

import (
	"context"
	"time"

	"banking-api/internal/repository"
	"banking-api/pkg/logger"
)

func StartCreditPaymentScheduler(repo *repository.CreditRepository, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			logger.Info("Running credit payment scheduler...")
			err := repo.CheckDuePayments(context.Background(), time.Now())
			if err != nil {
				logger.Error("Scheduler error: ", err)
			}
		}
	}()
}
