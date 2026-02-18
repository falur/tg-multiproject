package bot

import (
	"fmt"

	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleModePlan(c tele.Context) error {
	return b.setMode(c, "plan")
}

func (b *Bot) handleModeEdit(c tele.Context) error {
	return b.setMode(c, "edit")
}

func (b *Bot) setMode(c tele.Context, mode string) error {
	us := b.state.Get(c.Sender().ID)
	if us.ActiveProject == nil {
		return c.Respond(&tele.CallbackResponse{Text: "Проект не выбран"})
	}

	if err := b.store.UpdateProjectMode(*us.ActiveProject, mode); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка: " + err.Error()})
	}

	us.Mode = mode

	project, err := b.store.GetProject(*us.ActiveProject)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка: " + err.Error()})
	}

	_ = c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("Режим: %s", mode)})

	text := fmt.Sprintf("Проект: *%s*\nРежим: *%s*\nПуть: `%s`\n\nОтправьте задачу или используйте кнопки ниже.",
		project.Name, project.Mode, project.Path)
	return c.Edit(text, projectContextKeyboard(mode), tele.ModeMarkdown)
}

func (b *Bot) handleText(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)

	switch us.Step {
	case state.StepCreateName, state.StepCreateURL:
		return b.handleProjectTextInput(c, us)
	case state.StepInProject:
		return b.handleTaskSubmission(c, us)
	case state.StepRunning:
		return c.Send("Задача уже выполняется. Используйте /cancel для отмены.")
	default:
		return c.Send("Используйте /start для начала.", mainMenuKeyboard())
	}
}
