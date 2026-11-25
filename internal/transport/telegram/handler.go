package telegram

import (
	"Reminder_TelegramBot/internal/domain"
	"Reminder_TelegramBot/internal/service"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ReminderHandler struct {
	service *service.ReminderService
	bot     *bot.Bot
}

func NewReminderHandler(s *service.ReminderService) *ReminderHandler {
	return &ReminderHandler{
		service: s,
	}
}

func (h *ReminderHandler) HandleNewRemind(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	messageParts := strings.SplitN(update.Message.Text, " ", 3)

	if len(messageParts) < 3 {
		SendBotMsg(ctx, b, update, "Формат сообщения должен быть таким: </действие> <интервал> <текст>")
		return
	}

	delay, err := time.ParseDuration(messageParts[1])
	reminderTime := time.Now().Add(delay)

	if err != nil {
		SendBotMsg(ctx, b, update, "Неверно указано время. Используй формат вроде 10s, 15m или 1h")
		return
	}

	reminder, err := h.service.CreateReminder(ctx, update.Message.Chat.ID, messageParts[2], reminderTime)

	if errors.Is(err, domain.ErrNullText) {
		SendBotMsg(ctx, b, update, "Пустой текст напоминания")
		return
	}
	if errors.Is(err, domain.ErrReminderAtPast) {
		SendBotMsg(ctx, b, update, "Время напоминания установлено в прошлом")
		return
	}
	if err != nil {
		SendBotMsg(ctx, b, update, err.Error())
		return
	}

	SendBotMsg(ctx, b, update, "установлено напоминание "+reminder.Text+" на время "+reminderTime.Format("2006-01-02 15:04:05"))

}

func (h *ReminderHandler) HandleGetList(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	remindersList, err := h.service.GetPendingReminders(ctx)
	if err != nil {
		SendBotMsg(ctx, b, update, err.Error())
		return
	}
	if len(remindersList) == 0 {
		SendBotMsg(ctx, b, update, "Нет активных напоминаний")
		return
	}
	for _, reminder := range remindersList {
		SendBotMsg(ctx, b, update, reminder.Text+" на время "+reminder.Time.Format("2006-01-02 15:04:05"))
	}

}

func (h *ReminderHandler) HandleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	SendBotMsg(ctx, b, update, "Неизвестная команда. Используй /remind или /list»")
}

func SendBotMsg(ctx context.Context, b *bot.Bot, update *models.Update, message string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
}

func (h *ReminderHandler) SendNotification(ctx context.Context, chatId int64, message string) error {
	_, err := h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      "<b>(!)НАПОМИНАНИЕ: </b>" + message,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *ReminderHandler) SetBot(b *bot.Bot) {
	h.bot = b
}

func (h *ReminderHandler) RegisterCommands(b *bot.Bot) {
	//RegisterCommands commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remind", bot.MatchTypePrefix, h.HandleNewRemind)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypePrefix, h.HandleGetList)
}
