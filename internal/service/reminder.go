package service

import (
	"Reminder_TelegramBot/internal/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder domain.Reminder) (domain.Reminder, error)
	GetByID(ctx context.Context, id int64) (domain.Reminder, error)
	GetPending(ctx context.Context) ([]domain.Reminder, error)
	MarkAsSent(ctx context.Context, id int64) error
}

type ReminderService struct {
	repo ReminderRepository
}

func NewReminderService(repo ReminderRepository) *ReminderService {
	return &ReminderService{repo: repo}
}

func (r *ReminderService) CreateReminder(ctx context.Context, chatID int64, text string, notifyAt time.Time) (domain.Reminder, error) {
	if strings.TrimSpace(text) == "" {
		return domain.Reminder{}, domain.ErrNullText
	}
	if notifyAt.Before(time.Now()) {
		return domain.Reminder{}, domain.ErrReminderAtPast
	}
	reminder := domain.Reminder{
		ChatID: chatID,
		Text:   text,
		Time:   notifyAt,
		Status: domain.StatusPending,
	}
	reminder, err := r.repo.Create(ctx, reminder)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("create reminder error: %w", err)
	}
	return reminder, nil
}

func (r *ReminderService) GetPendingReminders(ctx context.Context) ([]domain.Reminder, error) {
	pendingReminders, err := r.repo.GetPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("there was an error getting pending reminders: %w", err)
	}
	return pendingReminders, nil
}
