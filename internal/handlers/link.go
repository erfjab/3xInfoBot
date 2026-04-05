package handlers

import (
	"fmt"
	"strings"
	"time"

	"3xinfobot/internal/panel"

	"github.com/erfjab/egobot/core"
	"github.com/erfjab/egobot/models"
)

func LinkHandler(b *core.Bot, update *models.Update, ctx *core.Context) error {
	text := update.Message.Text
	if !verifyLink(text) {
		_, err := b.SendMessage(&models.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please send a valid link starting with <code>vless://</code>",
			ParseMode: "HTML",
		})
		return err
	}

	uuid, err := extractUUID(text)
	if err != nil {
		_, err = b.SendMessage(&models.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "Could not extract UUID from the link.",
			ParseMode: "HTML",
		})
		return err
	}

	trafficResp, err := panel.DefaultClient.GetClientTrafficByUUID(uuid)
	if err != nil || !trafficResp.Success || len(trafficResp.Obj) == 0 {
		_, err = b.SendMessage(&models.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "No data found for this UUID.",
			ParseMode: "HTML",
		})
		return err
	}

	t := trafficResp.Obj[0]

	status := clientStatus(t)

	totalStr := "Unlimited"
	if t.Total > 0 {
		totalStr = formatBytes(t.Total)
	}

	remainingStr := "Unlimited"
	if t.Total > 0 {
		rem := t.Total - (t.Up + t.Down)
		if rem < 0 {
			rem = 0
		}
		remainingStr = formatBytes(rem)
	}

	expiryStr := "Never"
	if t.ExpiryTime > 0 {
		expiryStr = time.UnixMilli(t.ExpiryTime).Format("2006-01-02")
	}

	msg := fmt.Sprintf(
		"<b>❕ Status:</b> <code>%s</code>\n"+
			"<b>• Name:</b> <code>%s</code>\n"+
			"<b>• Used:</b> <code>%s</code>\n"+
			"<b>• Limit:</b> <code>%s</code>\n"+
			"<b>• Remaining:</b> <code>%s</code>\n"+
			"<b>• Expiry:</b> <code>%s</code>",
		status,
		t.Email,
		formatBytes(t.Up+t.Down),
		totalStr,
		remainingStr,
		expiryStr,
	)

	_, err = b.SendMessage(&models.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: "HTML",
	})
	return err
}

func clientStatus(t panel.ClientTraffic) string {
	if !t.Enable {
		return "Disabled"
	}
	if t.ExpiryTime > 0 && time.Now().UnixMilli() > t.ExpiryTime {
		return "Expired"
	}
	if t.Total > 0 && (t.Up+t.Down) >= t.Total {
		return "Limited"
	}
	return "Active"
}

func verifyLink(link string) bool {
	return strings.HasPrefix(link, "vless://")
}

func extractUUID(link string) (string, error) {
	withoutScheme := strings.TrimPrefix(link, "vless://")
	atIndex := strings.Index(withoutScheme, "@")
	if atIndex <= 0 {
		return "", fmt.Errorf("no UUID found")
	}
	return withoutScheme[:atIndex], nil
}

func formatBytes(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/KB)
	default:
		return fmt.Sprintf("%d B", b)
	}
}