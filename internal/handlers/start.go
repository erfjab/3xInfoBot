package handlers

import (
	"github.com/erfjab/egobot/core"
	"github.com/erfjab/egobot/models"
)

func StartHandler(b *core.Bot, update *models.Update, ctx *core.Context) error {
	_, err := b.SendMessage(&models.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Hello! Please send your link to get information about it.",
	})
	return err
}