package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"banking-api/internal/model"
	"banking-api/internal/repository"
	"banking-api/internal/security"
	"banking-api/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type CardService struct {
	cardRepo    *repository.CardRepository
	accountRepo *repository.AccountRepository
}

func NewCardervice(cardRepo *repository.CardRepository, accountRepo *repository.AccountRepository) *CardService {
	return &CardService{cardRepo: cardRepo, accountRepo: accountRepo}
}

func (s *CardService) CreateCard(ctx context.Context, userID, accountID int64) (*model.CardPlain, error) {

	// check if account belongs to user
	ownerID, err := s.accountRepo.GetOwnerID(ctx, accountID)
	if err != nil || ownerID != userID {
		return nil, errors.New("access denied")
	}

	// generate card payload
	number := utils.GenerateCardNumber("4276") // Visa test
	now := time.Now()
	threeYearsAfter := now.AddDate(3, 0, 0)
	expiry := threeYearsAfter.Format("1/06")
	cvv := fmt.Sprintf("%03d", rand.Intn(1000))

	// encrypt
	numberEnc, _ := security.EncryptPGP(number)
	_ = os.WriteFile("test_card_number.asc", numberEnc, 0644)
	expiryEnc, _ := security.EncryptPGP(expiry)
	cvvHash, _ := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	hmac, _ := security.GenerateHMAC(number + expiry + cvv)

	// save
	card, err := s.cardRepo.Create(ctx, userID, accountID, numberEnc, expiryEnc, string(cvvHash), hmac)
	if err != nil {
		return nil, err
	}

	return &model.CardPlain{
		ID:        card.ID,
		AccountID: accountID,
		Number:    number,
		Expiry:    expiry,
		CVV:       cvv,
	}, nil
}

func (s *CardService) GetUserCards(ctx context.Context, userID int64) ([]*model.CardRaw, error) {
	return s.cardRepo.GetByUser(ctx, userID)
}
