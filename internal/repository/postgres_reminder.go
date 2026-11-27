package repository

import (
	"Reminder_TelegramBot/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresReminderRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresReminderRepo(pool *pgxpool.Pool) *PostgresReminderRepo {
	return &PostgresReminderRepo{pool: pool}
}

func (r *PostgresReminderRepo) Create(ctx context.Context, reminder domain.Reminder) (domain.Reminder, error) {
	// 1. Получаем внутренний ID клиента или создаём его (UPSERT)
	var userID int64
	userQuery := `
		INSERT INTO users (telegram_chat_id) VALUES ($1)
		ON CONFLICT (telegram_chat_id) DO UPDATE SET telegram_chat_id = EXCLUDED.telegram_chat_id
		RETURNING id`

	err := r.pool.QueryRow(ctx, userQuery, reminder.ChatID).Scan(&userID)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("upsert user: %w", err)
	}

	// 2. Вставляем заказ и сразу забираем присвоенный базой ID
	reminderQuery := `
		INSERT INTO reminders (user_id, text, notify_at, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err = r.pool.QueryRow(ctx, reminderQuery, userID, reminder.Text, reminder.Time, reminder.Status).Scan(&reminder.ID)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("insert reminder: %w", err)
	}

	return reminder, nil
}

func (r *PostgresReminderRepo) Delete(ctx context.Context, id, chatID int64) error {
	query := `
		DELETE FROM reminders r
		USING users u
		WHERE r.user_id = u.id
		  AND r.id = $1
		  AND u.telegram_chat_id = $2`
	cmdTag, err := r.pool.Exec(ctx, query, id, chatID)

	if err != nil {
		return fmt.Errorf("error on deleting reminder : %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresReminderRepo) GetByID(ctx context.Context, id int64) (domain.Reminder, error) {

	var reminder domain.Reminder
	query := `
		SELECT r.id,u.telegram_chat_id, r.text, r.notify_at, r.status
			FROM reminders r 
			JOIN users u ON r.user_id = u.id
			WHERE r.id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&reminder.ID, &reminder.ChatID, &reminder.Text, &reminder.Time, &reminder.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Reminder{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("error getting remind by ID: %w", err)
	}

	return reminder, nil
}

func (r *PostgresReminderRepo) GetPending(ctx context.Context) ([]domain.Reminder, error) {

	var reminders []domain.Reminder
	query := `
		SELECT r.id,u.telegram_chat_id, r.text, r.notify_at, r.status
			FROM reminders r 
			JOIN users u ON r.user_id = u.id
			WHERE r.status = $1
			ORDER BY r.notify_at ASC`
	rows, err := r.pool.Query(ctx, query, domain.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("get pending reminder: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reminder domain.Reminder
		err = rows.Scan(&reminder.ID, &reminder.ChatID, &reminder.Text, &reminder.Time, &reminder.Status)
		if err != nil {
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		reminders = append(reminders, reminder)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("getting rows from db: %w", err)
	}

	return reminders, nil
}

func (r *PostgresReminderRepo) MarkAsSent(ctx context.Context, id int64) error {
	cmdTag, err := r.pool.Exec(ctx, "UPDATE reminders SET status = $1 WHERE id = $2", domain.StatusSent, id)

	if err != nil {
		return fmt.Errorf("error on marking reminder as SENT: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
