package repositories

import (
	"GameWala-Arcade/utils"
	"database/sql"
	"fmt"
)

type ArcadeRepository interface {
	ValidateArcade(arcadeId string) (bool, error)
	GetRazorpayAccountID(arcadeId string) (string, error)
}

type arcadeRepository struct {
	db *sql.DB
}

func NewArcadeRepository(db *sql.DB) *arcadeRepository {
	return &arcadeRepository{db: db}
}

func (r *arcadeRepository) ValidateArcade(arcadeId string) (bool, error) {
	// Implement the logic to validate the arcade ID against the database.
	// For now, we will return true for any non-empty arcade ID.
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM "Arcade" WHERE id = $1`, arcadeId).Scan(&count)
	if err != nil {
		utils.LogError("Getting some error while executing statement on DB %v", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *arcadeRepository) GetRazorpayAccountID(arcadeId string) (string, error) {
	var razorpayAccountID sql.NullString // Handles potential NULL values in database safely

	query := `SELECT get_razorpay_account_id($1);`

	// Assuming s.db is your *sql.DB instance
	err := r.db.QueryRow(query, arcadeId).Scan(&razorpayAccountID)
	if err != nil {
		return "", fmt.Errorf("failed to execute get_razorpay_account_id function: %w", err)
	}

	// Check if the arcade exists or if the client has a configured ID
	if !razorpayAccountID.Valid || razorpayAccountID.String == "" {
		return "", fmt.Errorf("razorpay account ID is missing or arcade ID '%s' is invalid", arcadeId)
	}

	return razorpayAccountID.String, nil
}
