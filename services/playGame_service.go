package services

import (
	"GameWala-Arcade/models"
	"GameWala-Arcade/repositories"
	"GameWala-Arcade/utils"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var maxTimeForLevelBoundedGame = uint16(120)

const staticStartingCode = "ABXYSO"

type PlayGameService interface {
	GetGames(arcadeId string) ([]models.GameResponse, error)
	// CheckGameCode(code string) (models.GameDetails, error) // arcade will hit this api
	// GenerateCode() (string, error)
}

type playGameService struct {
	playGameRepository repositories.PlayGameRepository
	redisClient        *redis.Client
}

func NewPlayGameService(playGameRepository repositories.PlayGameRepository,
	redisClient *redis.Client) *playGameService {
	return &playGameService{playGameRepository: playGameRepository, redisClient: redisClient}
}

func (s *playGameService) GetGames(arcadeId string) ([]models.GameResponse, error) {
	utils.LogInfo("Fetching all games from service for arcade ID: %s", arcadeId)
	games, err := s.playGameRepository.GetGames(arcadeId)

	if err != nil {
		utils.LogError("Failed to fetch games: %v", err)
		return nil, err
	}

	prices, err := s.playGameRepository.FetchPrices()

	for game := 0; game < len(games); game++ {
		currId := games[game].GameId
		if len(prices.TimeMap[currId]) > 0 {
			games[game].Price.ByTime = append(games[game].Price.ByTime, prices.TimeMap[currId]...)
		} else if len(prices.LevelMap[currId]) > 0 {
			games[game].Price.ByLevel = append(games[game].Price.ByLevel, prices.LevelMap[currId]...)
		}
	}

	utils.LogInfo("Successfully fetched %d games", len(games))
	return games, nil
}

func (s *playGameService) CheckGameCode(code string) (models.GameDetails, error) {
	if code == "" {
		utils.LogError("empty code in service layer? something's fishy 🐠")
		return models.GameDetails{}, fmt.Errorf("Code is empty")
	}

	status, err := s.playGameRepository.CheckGameCode(code)

	if err != nil {
		utils.LogError("Something went wrong... hmm BL layer, kinda sv issue?, err: %s", err)
		return status, err
	}

	return status, err
}

// func (s *playGameService) GenerateCode() (string, error) {
// 	ctx := context.Background()
// 	latestCode, err := s.redisClient.Get(ctx, "latest_arcade_code").Result()

// 	if err == redis.Nil {
// 		latestCode = staticStartingCode                             // starting code
// 		s.redisClient.Set(ctx, "latest_arcade_code", latestCode, 0) // 0 for no expiration
// 		return latestCode, err
// 	}

// 	newCode := getNextConsecutiveCode(latestCode)
// 	s.redisClient.Set(ctx, "latest_arcade_code", newCode, 0) // 0 for no expiration
// 	return newCode, nil
// }

// func getNextConsecutiveCode(code string) string {
// 	charset := []rune{'A', 'B', 'O', 'S', 'X', 'Y'}
// 	base := len(charset)
// 	runes := []rune(code)
// 	n := len(runes)

// 	carry := 1
// 	for i := n - 1; i >= 0; i-- {
// 		if carry == 0 {
// 			break
// 		}
// 		idx := indexOf(charset, runes[i])
// 		if idx == -1 {
// 			idx = 0
// 		}
// 		newIdx := (idx + carry) % base
// 		carry = (idx + carry) / base
// 		runes[i] = charset[newIdx]
// 	}
// 	return string(runes)
// }

// func indexOf(slice []rune, r rune) int {
// 	for i, v := range slice {
// 		if v == r {
// 			return i
// 		}
// 	}
// 	return -1
// }
