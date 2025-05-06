package service

import (
	"banking-api/internal/model"
	"banking-api/internal/repository"
	"context"
	"errors"
	"time"
)

type AnalyticsService struct {
	repo *repository.AnalyticsRepository
}

func NewAnalyticsService(repo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetMonthlyStats(ctx context.Context, userID int64) (*model.MonthlyStats, error) {
	stats, err := s.repo.GetMonthlyStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	dueTotal, err := s.repo.GetUnpaidCreditTotal(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.DueTotal = dueTotal
	return stats, nil
}

func (s *AnalyticsService) GetPredictedBalance(ctx context.Context, userID int64, days int) (float64, error) {
	if days < 0 || days > 365 {
		return 0, errors.New("prediction period must be between 0 and 365")
	}
	until := time.Now().AddDate(0, 0, days)
	return s.repo.GetPredictedBalance(ctx, userID, until)
}
