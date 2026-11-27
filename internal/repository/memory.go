package repository

import (
	"Reminder_TelegramBot/internal/domain"
	"context"
	"sync"
)

type InMemoryReminderRepo struct {
	mu        sync.RWMutex
	reminders map[int64]domain.Reminder
	nextID    int64
}

func NewInMemoryReminderRepo() *InMemoryReminderRepo {
	return &InMemoryReminderRepo{
		reminders: make(map[int64]domain.Reminder),
		nextID:    1,
	}

}

func (r *InMemoryReminderRepo) Create(ctx context.Context, reminder domain.Reminder) (domain.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	reminder.ID = r.nextID
	r.nextID++
	r.reminders[reminder.ID] = reminder

	return reminder, nil
}

func (r *InMemoryReminderRepo) Delete(ctx context.Context, id, ChatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if rem, ok := r.reminders[id]; !ok || rem.ChatID != ChatID {
		return domain.ErrNotFound
	}
	delete(r.reminders, id)

	return nil
}

func (r *InMemoryReminderRepo) GetByID(ctx context.Context, id int64) (domain.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reminder, ok := r.reminders[id]
	if !ok {
		return domain.Reminder{}, domain.ErrNotFound
	}
	return reminder, nil
}

func (r *InMemoryReminderRepo) GetPending(ctx context.Context) ([]domain.Reminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reminders := make([]domain.Reminder, 0, len(r.reminders))

	for _, reminder := range r.reminders {
		if reminder.Status == domain.StatusPending {
			reminders = append(reminders, reminder)
		}
	}

	return reminders, nil
}

func (r *InMemoryReminderRepo) MarkAsSent(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	reminder, err := r.GetByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	reminder.Status = domain.StatusSent

	return nil
}
