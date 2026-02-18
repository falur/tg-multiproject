package bot

import tele "gopkg.in/telebot.v3"

func (b *Bot) handleStart(c tele.Context) error {
	b.state.Reset(c.Sender().ID)
	return c.Send("Добро пожаловать! Выберите действие:", mainMenuKeyboard())
}
