package services

import (
	"GameWala-Arcade/repositories"
)

type ArcadeService interface {
	ValidateArcade(arcadeId string) (bool, error)
	GetRazorpayAccountID(arcadeId string) (string, error)
}

type arcadeService struct {
	arcadeRepository repositories.ArcadeRepository
}

func NewArcadeService(arcadeRepository repositories.ArcadeRepository) *arcadeService {
	return &arcadeService{
		arcadeRepository: arcadeRepository,
	}
}

func (s *arcadeService) ValidateArcade(arcadeId string) (bool, error) {
	return s.arcadeRepository.ValidateArcade(arcadeId)
}

func (s *arcadeService) GetRazorpayAccountID(arcadeId string) (string, error) {
	razorpay_account_id, err := s.arcadeRepository.GetRazorpayAccountID(arcadeId)
	if err != nil {
		return razorpay_account_id, err
	}

	return razorpay_account_id, err
}
