package services

import (
	"GameWala-Arcade/models"
	"GameWala-Arcade/repositories"
	"GameWala-Arcade/utils"
	mqtt "GameWala-Arcade/utils/mqtt"
	"fmt"
)

type ConnectionToBrokerService interface {
	PublishMessage(arcadeId string, gameStatus models.GameStatus) error
}

type connectionToBrokerService struct {
	mqtt               *mqtt.MQTTService
	playGameRepository repositories.PlayGameRepository
}

func NewConnectionToBrokerService(mqttService *mqtt.MQTTService, playGameRepository repositories.PlayGameRepository) *connectionToBrokerService {
	return &connectionToBrokerService{mqtt: mqttService, playGameRepository: playGameRepository}
}

func (s *connectionToBrokerService) PublishMessage(arcadeId string, gameStatus models.GameStatus) error {

	// fetch game details on the basis of game id.
	payload, err := s.playGameRepository.FetchGameDetails(gameStatus.GameId)

	if err != nil {
		utils.LogError("Failed to fetch game details for game ID %d: %v", gameStatus.GameId, err)
		return fmt.Errorf("error fetching game details: %w", err)
	}

	gameDataMapper(&payload, gameStatus)

	// Implement the logic to publish a message to the specified topic on the broker
	var topic = fmt.Sprintf("arcade/%s/game/command", arcadeId)
	cmd := models.Command{Action: "launch", Data: payload}

	utils.LogInfo("Publishing message to topic.... '%s'\n", topic)

	err = s.mqtt.Publish(topic, cmd)
	if err != nil {
		utils.LogError("could not publish message to broker, error: %v", err)
		return err
	}

	return err
}

func gameDataMapper(payload *models.GameDetails, gameStatus models.GameStatus) {
	payload.IsPlayed = gameStatus.IsPlayed
	payload.Time = gameStatus.PlayTime
	payload.Level = gameStatus.Levels
	payload.IsTimed = gameStatus.IsTimed
}
