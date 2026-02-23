package services

import (
	"GameWala-Arcade/models"
	"GameWala-Arcade/repositories"
	"GameWala-Arcade/utils"
)

type HandlePaymentService interface {
	SaveOrderDetails(models.PaymentStatus) error
	SaveGameStatus(status models.GameStatus) (int, error)
}

type handlePaymentService struct {
	handlePaymentRepository repositories.HandlePaymentRepository
	playGameRepository      repositories.PlayGameRepository
}

func NewHandlePaymentService(handlePaymentRepository repositories.HandlePaymentRepository, playGameRepository repositories.PlayGameRepository) *handlePaymentService {
	return &handlePaymentService{handlePaymentRepository: handlePaymentRepository, playGameRepository: playGameRepository}
}

func (s *handlePaymentService) SaveOrderDetails(details models.PaymentStatus) error {
	err := s.handlePaymentRepository.SaveOrderDetails(details)
	return err
}

func (s *handlePaymentService) SaveGameStatus(status models.GameStatus) (int, error) {
	utils.LogInfo("Processing save game status for game ID %d", status.GameId)

	if status.IsTimed && status.PlayTime != nil {
		err := s.validateTimeAndPrice(status.GameId, status.Price, status.PlayTime)

		if err != nil {
			utils.LogError("Time and price validation failed for game ID %d: %v", status.GameId, err)
			return 2, err // 2 means, time and price didn't match (convert this to enum later)
		}
	} else {
		err := s.validateLevelsAndPrice(status.GameId, status.Price, status.Levels)

		if err != nil {
			utils.LogError("Level and price validation failed for game ID %d: %v", status.GameId, err)
			return 3, err // 3 means, level and price didn't match (convert this to enum later)
		}
	}

	return s.handlePaymentRepository.SaveGameStatus(status)
}

func (s *handlePaymentService) validateTimeAndPrice(gameId uint16, price uint16, playTime *uint16) error {
	utils.LogInfo("Validating time and price for game ID %d: price=%d, time=%d", gameId, price, *playTime)
	//call db to cheeck if time and price match with the feeded value.
	err := s.playGameRepository.ValidateTimeAndPrice(gameId, price, playTime)

	return err
}

func (s *handlePaymentService) validateLevelsAndPrice(gameId uint16, price uint16, levels *uint8) error {
	utils.LogInfo("Validating levels and price for game ID %d: price=%d, level=%d", gameId, price, *levels)
	//call db to cheeck if level and price match with the feeded value.
	err := s.playGameRepository.ValidateLevelsAndPrice(gameId, price, levels)

	return err
}
