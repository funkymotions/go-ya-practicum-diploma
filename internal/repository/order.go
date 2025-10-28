package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	apperrors "github.com/funkymotions/go-ya-practicum-diploma/internal/apperror"
	"github.com/funkymotions/go-ya-practicum-diploma/internal/model"
)

type orderRepository struct {
	driver SQLExecutor
}

func NewOrderRepository(d SQLExecutor) *orderRepository {
	return &orderRepository{
		driver: d,
	}
}

func (r *orderRepository) FindOneByID(orderID string) (*model.Order, error) {
	var order model.Order
	sqlString := `
		SELECT
			id,
			user_id,
			accrual,
			order_status,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			id = $1;`

	row := r.driver.QueryRow(sqlString, orderID)
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Accrual,
		&order.OrderStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindManyByStatus(statuses ...model.OrderStatus) (*[]model.Order, error) {
	if len(statuses) == 0 {
		return nil, errors.New("no statuses provided")
	}
	statusInClause := make([]string, 0, len(statuses))
	statusArgs := make([]interface{}, len(statuses))
	for i, s := range statuses {
		statusInClause = append(statusInClause, fmt.Sprintf("$%d", i+1))
		statusArgs[i] = s
	}
	inClause := strings.Join(statusInClause, ", ")
	sqlString := `
		SELECT
			id,
			user_id,
			accrual,
			order_status,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			order_status IN (` + inClause + `);`
	orders := make([]model.Order, 0)
	rows, err := r.driver.Query(sqlString, statusArgs...)
	if err == sql.ErrNoRows {
		return &orders, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Accrual,
			&order.OrderStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &orders, nil
}

func (r *orderRepository) CreateOrder(orderID string, userID uint) (*model.Order, error) {
	var order model.Order
	sqlString := `
		INSERT INTO
			orders (
				id,
				user_id,
				order_status
			)
		VALUES
			($1, $2, $3)
		ON CONFLICT
			(id)
		DO NOTHING
		RETURNING
			id,
			user_id,
			accrual,
			order_status,
			created_at,
			updated_at;`

	row := r.driver.QueryRow(
		sqlString,
		orderID,
		userID,
		model.StatusNew,
	)
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Accrual,
		&order.OrderStatus,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrOrderAlreadyRegistered
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetUserOrders(userID uint) (*[]model.Order, error) {
	var orders []model.Order
	sqlString := `
		SELECT
			id,
			user_id,
			accrual,
			order_status,
			created_at,
			updated_at
		FROM
			orders
		WHERE
			user_id = $1
		ORDER BY
			created_at ASC;`

	rows, err := r.driver.Query(sqlString, userID)
	if err == sql.ErrNoRows {
		return &[]model.Order{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Accrual,
			&order.OrderStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &orders, nil
}

func (r *orderRepository) GetUserOrderBalance(userID uint) (float64, error) {
	sqlString := `
		SELECT
			COALESCE(SUM(accrual), 0) as balance
		FROM
			orders
		WHERE
			user_id = $1
		AND
			order_status = $2;`
	var balance float64
	err := r.driver.QueryRow(sqlString, userID, model.StatusProcessed).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *orderRepository) UpdateOrder(order *model.Order) (*model.Order, error) {
	sqlString := `
		UPDATE
			orders
		SET
			accrual = $1,
			order_status = $2,
			updated_at = NOW()
		WHERE
			id = $3
		RETURNING
			id,
			user_id,
			accrual,
			order_status,
			created_at,
			updated_at;`
	row := r.driver.QueryRow(
		sqlString,
		order.Accrual,
		order.OrderStatus,
		order.ID,
	)
	var updatedOrder model.Order
	err := row.Scan(
		&updatedOrder.ID,
		&updatedOrder.UserID,
		&updatedOrder.Accrual,
		&updatedOrder.OrderStatus,
		&updatedOrder.CreatedAt,
		&updatedOrder.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updatedOrder, nil
}
