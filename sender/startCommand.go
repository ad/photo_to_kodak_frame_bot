package sender

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Sender) HandleStartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("/start message from", update.Message.Chat.ID)

	go func(ctx context.Context, b *bot.Bot, update *models.Update) {
		s.parseStartCommand(update)
	}(ctx, b, update)
}

func (s *Sender) parseStartCommand(update *models.Update) error {
	if update.Message == nil {
		return nil
	}

	if len(s.config.TelegramAdminIDsList) != 0 {
		if update.Message.From != nil {
			if !slices.Contains(s.config.TelegramAdminIDsList, update.Message.From.ID) {
				return nil
			}
		}
	}

	s.MakeRequestDeferred(DeferredMessage{
		Method: "sendMessage",
		ChatID: update.Message.From.ID,
		Text:   "hello, send me a photo",
	}, s.SendResult)

	return nil
}
