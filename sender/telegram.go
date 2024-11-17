package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/ad/photo_to_kodak_frame_bot/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Sender struct {
	sync.RWMutex
	logger           *slog.Logger
	config           *config.Config
	Bot              *bot.Bot
	BotUser          *models.User
	deferredMessages map[int64]chan DeferredMessage
	lastMessageTimes map[int64]int64
}

func InitSender(ctx context.Context, logger *slog.Logger, config *config.Config) (*Sender, error) {
	sender := &Sender{
		logger:           logger,
		config:           config,
		deferredMessages: make(map[int64]chan DeferredMessage),
		lastMessageTimes: make(map[int64]int64),
	}

	opts := []bot.Option{
		bot.WithSkipGetMe(),
		bot.WithAllowedUpdates([]string{"callback_query", "message", "chat_member", "chat_join_request"}),
		bot.WithDefaultHandler(sender.handler),
	}

	b, newBotError := bot.New(config.TelegramToken, opts...)
	if newBotError != nil {
		return nil, fmt.Errorf("start bot error: %s", newBotError)
	}

	me, err := b.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	sender.BotUser = me

	go b.Start(ctx)
	go sender.sendDeferredMessages()

	sender.Bot = b

	return sender, nil
}

func (s *Sender) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if s.config.Debug {
		jsonStr, _ := json.Marshal(update)
		fmt.Println("update", string(jsonStr))

		s.logger.Debug(formatUpdateForLog(update))
	}

	if len(s.config.TelegramAdminIDsList) != 0 {
		if update.Message.From != nil {
			if !slices.Contains(s.config.TelegramAdminIDsList, update.Message.From.ID) {
				return
			}
		}
	}

	if update.Message != nil && update.Message.Photo != nil {
		go func(ctx context.Context, b *bot.Bot, update *models.Update) {
			s.processPhotos(ctx, b, update)
		}(ctx, b, update)

		return
	}
}
