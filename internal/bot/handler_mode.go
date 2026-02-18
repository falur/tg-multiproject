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
		return c.Respond(&tele.CallbackResponse{Text: "No project selected"})
	}

	if err := b.store.UpdateProjectMode(*us.ActiveProject, mode); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Error: " + err.Error()})
	}

	us.Mode = mode

	project, err := b.store.GetProject(*us.ActiveProject)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Error: " + err.Error()})
	}

	_ = c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("Mode: %s", mode)})

	text := fmt.Sprintf("Project: *%s*\nMode: *%s*\nPath: `%s`\n\nSend a task or use the buttons below.",
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
		return c.Send("A task is already running. Use /cancel to stop it.")
	default:
		return c.Send("Use /start to begin.", mainMenuKeyboard())
	}
}
