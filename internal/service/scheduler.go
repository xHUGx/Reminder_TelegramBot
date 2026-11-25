package service

import (
	"context"
	"log"
	"time"
)

type Notifier interface {
	SendNotification(ctx context.Context, chatId int64, message string) error
}
type ReminderScheduler struct {
	repo     ReminderRepository
	notifier Notifier
	interval time.Duration
}

func NewReminderScheduler(repo ReminderRepository, notifier Notifier, interval time.Duration) *ReminderScheduler {
	return &ReminderScheduler{
		repo:     repo,
		notifier: notifier,
		interval: interval,
	}
}

func (s *ReminderScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Планировщик уведомлений остановлен")
			return
		case <-ticker.C:
			s.checkAndSend(ctx)

		}
	}
}

func (s *ReminderScheduler) checkAndSend(ctx context.Context) {
	pendingReminders, err := s.repo.GetPending(ctx)
	if err != nil {
		log.Println("ошибка получения активных напоминаний:", err)
		return
	}
	now := time.Now()

	for _, reminder := range pendingReminders {
		if !now.Before(reminder.Time) {
			err := s.notifier.SendNotification(ctx, reminder.ChatID, reminder.Text)
			if err != nil {
				log.Println("ошибка отправки напоминания:", err)
				continue
			}
			err = s.repo.MarkAsSent(ctx, reminder.ID)
			if err != nil {
				log.Println("ошибка отметки сообщения как выполненого:", err)
				continue
			}
		} else {
			continue
		}
	}
}
