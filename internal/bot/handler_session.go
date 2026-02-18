package bot

import (
	"strconv"

	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleSessions(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)
	if us.ActiveProject == nil {
		return c.Respond(&tele.CallbackResponse{Text: "No project selected"})
	}

	sessions, err := b.store.ListSessions(*us.ActiveProject)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Error: " + err.Error()})
	}

	if len(sessions) == 0 {
		return c.Edit("No sessions yet.", projectContextKeyboard(us.Mode))
	}

	return c.Edit("Select a session to resume:", sessionListKeyboard(sessions))
}

func (b *Bot) handleResumeSession(c tele.Context) error {
	data := c.Callback().Data
	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Invalid session ID"})
	}

	sess, err := b.store.GetSession(id)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Session not found"})
	}

	us := b.state.Get(c.Sender().ID)
	us.SessionID = sess.SessionID
	us.Step = state.StepInProject

	_ = c.Respond(&tele.CallbackResponse{Text: "Session loaded"})
	return c.Edit("Session resumed. Send your next task.", projectContextKeyboard(us.Mode))
}
