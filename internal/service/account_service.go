package service

import (
	"banking-api/internal/model"
	"banking-api/internal/repository"
	"context"
	"errors"
	"fmt"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64) (*model.Account, error) {
	return s.accountRepo.Create(ctx, userID)
}

func (s *AccountService) Transfer(ctx context.Context, fromUserID, fromAccountID, toAccountID int64, amount float64) error {
	ownerID, err := s.accountRepo.GetOwnerID(ctx, fromAccountID)
	if err != nil {
		return fmt.Errorf("from account not found")
	}
	if ownerID != fromUserID {
		return errors.New("access denied")
	}
	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to same account")
	}
	if amount <= 0 {
		return errors.New("invalid transfer amount")
	}

	return s.accountRepo.TransferTx(ctx, fromAccountID, toAccountID, amount)
}

func (s *AccountService) Deposit(ctx context.Context, userID, accountID int64, amount float64) (*model.Account, error) {
	ownerID, err := s.accountRepo.GetOwnerID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if ownerID != userID {
		return nil, errors.New("access denied")
	}
	return s.accountRepo.Deposit(ctx, amount, accountID)
}

func (s *AccountService) Withdraw(ctx context.Context, userID, accountID int64, amount float64) (*model.Account, error) {
	ownerID, err := s.accountRepo.GetOwnerID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if ownerID != userID {
		return nil, errors.New("access denied")
	}
	return s.accountRepo.Withdraw(ctx, amount, accountID)
}

func (s *AccountService) Get(ctx context.Context, userID, accountID int64) (*model.Account, error) {
	ownerID, err := s.accountRepo.GetOwnerID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if ownerID != userID {
		return nil, errors.New("access denied")
	}
	return s.accountRepo.Get(ctx, accountID, userID)
}
