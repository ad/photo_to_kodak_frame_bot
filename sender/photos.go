package sender

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/go-pkgz/email"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func (s *Sender) processPhotos(ctx context.Context, b *bot.Bot, update *models.Update) error {
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

	fileID := ""
	size := 0
	for _, photoItem := range update.Message.Photo {
		if photoItem.FileSize > size {
			fileID = photoItem.FileID
		}
	}

	if fileID == "" {
		return nil
	}

	photoPath := getFilePath(ctx, b, fileID)

	if err := s.sendEmail(photoPath); err == nil {
		s.MakeRequestDeferred(DeferredMessage{
			Method: "sendMessage",
			ChatID: update.Message.From.ID,
			Text:   "Photo sent",
		}, s.SendResult)
	} else {
		s.MakeRequestDeferred(DeferredMessage{
			Method: "sendMessage",
			ChatID: update.Message.From.ID,
			Text:   fmt.Sprintf("Error sending email: %s", err.Error()),
		}, s.SendResult)
	}

	return nil
}

func getFilePath(ctx context.Context, b *bot.Bot, fileID string) string {
	fileModel, errGetFile := b.GetFile(ctx, &bot.GetFileParams{FileID: fileID})
	if errGetFile != nil {
		return ""
	}

	return fileModel.FilePath
}

func (s *Sender) sendEmail(photoPath string) error {
	filePath, err := downloadFile(fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", s.config.TelegramToken, photoPath))

	if err != nil {
		return err
	}

	client := email.NewSender(s.config.SMTP_HOST, email.Port(s.config.SMTP_PORT), email.STARTTLS(true), email.ContentType("text/plain"), email.Auth(s.config.SMTP_USER, s.config.SMTP_PASSWORD))
	err = client.Send("",
		email.Params{
			From:        s.config.FromEmail,
			To:          []string{s.config.TargetEmail},
			Subject:     "New photo",
			Attachments: []string{filePath},
		},
	)

	// remove temp file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error removing temp file: %s", err)
	}

	if err != nil {
		return fmt.Errorf("error sending email: %s", err)
	}

	return nil
}

func downloadFile(url string) (filepath string, err error) {
	// generate temp file name
	filepath = fmt.Sprintf("%s/%s", os.TempDir(), uuid.New().String())

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return filepath, err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return filepath, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return filepath, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return filepath, err
	}

	return filepath, nil
}
