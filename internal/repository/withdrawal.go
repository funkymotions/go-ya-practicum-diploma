package repository

import "github.com/funkymotions/go-ya-practicum-diploma/internal/model"

type withdrawalRepository struct {
	driver SQLExecutor
}

func NewWithdrawalRepository(d SQLExecutor) *withdrawalRepository {
	return &withdrawalRepository{
		driver: d,
	}
}

func (r *withdrawalRepository) CreateWithdrawal(userID uint, orderID string, amount float64) error {
	sqlString := `
		INSERT INTO
		withdrawals (
			user_id,
			order_id,
			amount
		)
		VALUES
			($1, $2, $3);`

	_, err := r.driver.Exec(sqlString, userID, orderID, amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *withdrawalRepository) GetUserWithdrawals(userID uint) (*[]model.Withdrawal, error) {
	sqlString := `
		SELECT
			id,
			user_id,
			order_id,
			amount,
			processed_at
		FROM
			withdrawals
		WHERE
			user_id = $1
		ORDER BY
			processed_at;`
	rows, err := r.driver.Query(sqlString, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.OrderID,
			&w.Amount,
			&w.ProcessedAt,
		); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &withdrawals, nil
}

func (r *withdrawalRepository) GetUserWithdrawalBalance(userID uint) (float64, error) {
	sqlString := `
		SELECT
			COALESCE(SUM(amount), 0) as balance
		FROM
			withdrawals
		WHERE
			user_id = $1;`

	row := r.driver.QueryRow(sqlString, userID)
	var balance float64
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}
