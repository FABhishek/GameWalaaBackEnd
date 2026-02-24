package repositories

import (
	"GameWala-Arcade/models"
	"GameWala-Arcade/utils"
	"database/sql"
	"fmt"
)

type HandlePaymentRepository interface {
	SaveOrderDetails(models.PaymentStatus) error
	SaveGameStatus(status models.GameStatus) (int, error)
}

type handlePaymentRepository struct {
	db *sql.DB
}

func NewHandlePaymentReposiory(db *sql.DB) *handlePaymentRepository {
	return &handlePaymentRepository{db: db}
}

func (r *handlePaymentRepository) SaveOrderDetails(details models.PaymentStatus) error {
	utils.LogInfo("Saving payment status for payment ID %s", details.RazorpayPaymentId)

	_, err := r.db.Exec("SELECT func_InsertPaymentStatus($1, $2, $3, $4)",
		details.OrderCreationId,
		details.RazorpayPaymentId,
		details.RazorpayOrderId,
		details.RazorpaySignature)

	if err != nil {
		utils.LogError("Failed to execute payment status for payment ID %s: %v", details.RazorpayPaymentId, err)
		return fmt.Errorf("error executing function: %w", err)
	}

	utils.LogInfo("Successfully saved payment status for order ID %s", details.OrderCreationId)
	return nil
}

func (r *handlePaymentRepository) SaveGameStatus(status models.GameStatus) (int, error) {
	utils.LogInfo("Saving game status to database for game ID %d", status.GameId)

	// Prepare the call to the stored procedure
	stmt, err := r.db.Prepare("SELECT func_InsertGameStatus($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		utils.LogError("Failed to prepare save game status statement: %v", err)
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status.GameId, status.Name, status.IsPlayed, status.Price,
		status.PlayTime, status.Levels, status.PaymentReference, status.ArcadeId)

	if err != nil {
		utils.LogError("Failed to execute save game status for game ID %d: %v", status.GameId, err)
		return 0, fmt.Errorf("error executing function: %w", err)
	}

	utils.LogInfo("Successfully saved game status for game ID %d", status.GameId)
	return 1, nil
}
