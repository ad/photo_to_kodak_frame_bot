package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/ad/photo_to_kodak_frame_bot/config"
	"github.com/ad/photo_to_kodak_frame_bot/logger"
	"github.com/ad/photo_to_kodak_frame_bot/sender"
	"github.com/go-telegram/bot"
)

func Run(ctx context.Context, w io.Writer, args []string) error {
	conf, errInitConfig := config.InitConfig(os.Args)
	if errInitConfig != nil {
		return errInitConfig
	}

	lgr := logger.InitLogger(conf.Debug)

	// Recovery
	defer func() {
		if p := recover(); p != nil {
			lgr.Error(fmt.Sprintf("panic recovered: %s; stack trace: %s", p, string(debug.Stack())))
		}
	}()

	sndr, errInitSender := sender.InitSender(ctx, lgr, conf)
	if errInitSender != nil {
		return errInitSender
	}

	if len(conf.TelegramAdminIDsList) != 0 {
		sndr.MakeRequestDeferred(sender.DeferredMessage{
			Method: "sendMessage",
			ChatID: conf.TelegramAdminIDsList[0],
			Text:   "Bot restarted: " + sndr.BotUser.Username,
		}, sndr.SendResult)
	}

	sndr.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, sndr.HandleStartCommand)

	return nil
}
