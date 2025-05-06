package service

import (
	"context"
	"errors"

	"banking-api/internal/model"
	"banking-api/internal/repository"
	"banking-api/pkg/utils"
)

type CreditService struct {
	accountRepo *repository.AccountRepository
	creditRepo  *repository.CreditRepository
}

func NewCreditService(accountRepo *repository.AccountRepository, creditRepo *repository.CreditRepository) *CreditService {
	return &CreditService{accountRepo: accountRepo, creditRepo: creditRepo}
}

func (s *CreditService) CreateCredit(ctx context.Context, userID, accountID int64, principal, annualRate float64, months int) (*model.Credit, error) {
	// check account owner
	ownerID, err := s.accountRepo.GetOwnerID(ctx, accountID)
	if err != nil || ownerID != userID {
		return nil, errors.New("access denied to account")
	}

	// calculate monthly payment
	payment := utils.CalculateAnnuity(principal, annualRate, months)

	// save credit
	credit := &model.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Principal:      principal,
		InterestRate:   annualRate,
		TermMonths:     months,
		MonthlyPayment: payment,
	}
	err = s.creditRepo.SaveCredit(ctx, credit)
	if err != nil {
		return nil, err
	}

	// Генерация графика
	err = s.creditRepo.GenerateSchedule(ctx, credit.ID, months, payment)
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *CreditService) GetSchedule(ctx context.Context, userID, creditID int64) ([]model.PaymentSchedule, error) {
	return s.creditRepo.GetScheduleByCreditID(ctx, userID, creditID)
}
