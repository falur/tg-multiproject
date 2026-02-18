package bot

import (
	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleCancel(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)
	if us.Step != state.StepRunning {
		return c.Send("Nothing is running.")
	}

	b.mu.Lock()
	cancel := b.cancelClaude
	b.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	return nil
}

func (b *Bot) handleCancelCallback(c tele.Context) error {
	_ = c.Respond(&tele.CallbackResponse{Text: "Cancelling..."})
	return b.handleCancel(c)
}
