package bot

import (
	"strconv"

	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleSessions(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)
	if us.ActiveProject == nil {
		return c.Respond(&tele.CallbackResponse{Text: "Проект не выбран"})
	}

	sessions, err := b.store.ListSessions(*us.ActiveProject)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка: " + err.Error()})
	}

	if len(sessions) == 0 {
		return c.Edit("Сессий пока нет.", projectContextKeyboard(us.Mode))
	}

	return c.Edit("Выберите сессию для продолжения:", sessionListKeyboard(sessions))
}

func (b *Bot) handleResumeSession(c tele.Context) error {
	data := c.Callback().Data
	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Неверный ID сессии"})
	}

	sess, err := b.store.GetSession(id)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Сессия не найдена"})
	}

	us := b.state.Get(c.Sender().ID)
	us.SessionID = sess.SessionID
	us.Step = state.StepInProject

	_ = c.Respond(&tele.CallbackResponse{Text: "Сессия загружена"})
	return c.Edit("Сессия восстановлена. Отправьте следующую задачу.", projectContextKeyboard(us.Mode))
}
