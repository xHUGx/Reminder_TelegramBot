package domain

import (
	"errors"
	"time"
)

type ReminderStatus string

const (
	StatusPending   ReminderStatus = "pending"
	StatusSent      ReminderStatus = "sent"
	StatusCancelled ReminderStatus = "cancelled"
)

type Reminder struct {
	ID     int64
	ChatID int64
	Text   string
	Time   time.Time
	Status ReminderStatus
}

var (
	ErrNotFound       = errors.New("reminder not found")
	ErrEmpty          = errors.New("reminder list is empty")
	ErrNullText       = errors.New("reminder list is null")
	ErrReminderAtPast = errors.New("reminder at past")
	ErrAlreadyExists  = errors.New("reminder already exists")
)
