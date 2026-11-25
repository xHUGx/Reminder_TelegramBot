package main

import (
	"Reminder_TelegramBot/internal/repository"
	"Reminder_TelegramBot/internal/service"
	"Reminder_TelegramBot/internal/transport/telegram"
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
)

// Send any text message to the bot after the bot has been started
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	//Initialization of services
	reminderRepo := repository.NewInMemoryReminderRepo()
	reminderService := service.NewReminderService(reminderRepo)
	reminderHandler := telegram.NewReminderHandler(reminderService)
	reminderScheduler := service.NewReminderScheduler(reminderRepo, reminderHandler, 1*time.Second)

	opts := []bot.Option{
		bot.WithDefaultHandler(reminderHandler.HandleDefault),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		log.Fatalf("ошибка запуска бота: %v", err)
	}

	reminderHandler.SetBot(b)
	reminderHandler.RegisterCommands(b)

	go reminderScheduler.Start(ctx)
	b.Start(ctx)
}
