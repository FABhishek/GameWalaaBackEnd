package services

import (
	"GameWala-Arcade/repositories"
)

type ArcadeService interface {
	ValidateArcade(arcadeId string) (bool, error)
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
