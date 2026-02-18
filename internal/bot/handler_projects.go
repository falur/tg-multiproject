package bot

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"tg-multiproject/internal/github"
	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleMyProjects(c tele.Context) error {
	projects, err := b.store.ListProjects()
	if err != nil {
		return c.Send("Error loading projects: " + err.Error())
	}
	if len(projects) == 0 {
		return c.Edit("No projects yet. Create one!", mainMenuKeyboard())
	}
	return c.Edit("Select a project:", projectListKeyboard(projects))
}

func (b *Bot) handleSelectProject(c tele.Context) error {
	data := c.Callback().Data
	id, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Invalid project ID"})
	}

	project, err := b.store.GetProject(id)
	if err != nil {
		return c.Respond(&tele.CallbackResponse{Text: "Project not found"})
	}

	us := b.state.Get(c.Sender().ID)
	us.Step = state.StepInProject
	us.ActiveProject = &project.ID
	us.Mode = project.Mode
	us.SessionID = ""

	text := fmt.Sprintf("Project: *%s*\nMode: *%s*\nPath: `%s`\n\nSend a task or use the buttons below.",
		project.Name, project.Mode, project.Path)

	return c.Edit(text, projectContextKeyboard(project.Mode), tele.ModeMarkdown)
}

func (b *Bot) handleCreateProject(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)
	us.Step = state.StepCreateName
	return c.Edit("Enter project name:")
}

func (b *Bot) handleSkip(c tele.Context) error {
	us := b.state.Get(c.Sender().ID)
	if us.Step != state.StepCreateURL {
		return nil
	}
	return b.finalizeProjectCreation(c, us, "")
}

func (b *Bot) handleBack(c tele.Context) error {
	b.state.Reset(c.Sender().ID)
	return c.Edit("Choose an action:", mainMenuKeyboard())
}

func (b *Bot) handleProjectTextInput(c tele.Context, us *state.UserState) error {
	switch us.Step {
	case state.StepCreateName:
		us.ProjectName = c.Text()
		us.Step = state.StepCreateURL
		return c.Send("Enter GitHub URL (or press Skip):", &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{{btnSkip}},
		})

	case state.StepCreateURL:
		return b.finalizeProjectCreation(c, us, c.Text())

	default:
		return nil
	}
}

func (b *Bot) finalizeProjectCreation(c tele.Context, us *state.UserState, githubURL string) error {
	name := us.ProjectName
	projectPath := filepath.Join(b.cfg.ProjectsDir, name)

	if githubURL != "" {
		if err := c.Send("Cloning repository..."); err != nil {
			return err
		}
		if err := github.Clone(githubURL, projectPath); err != nil {
			us.Step = state.StepCreateURL
			return c.Send("Clone failed: " + err.Error() + "\nTry again or press Skip.")
		}
	} else {
		if err := os.MkdirAll(projectPath, 0o755); err != nil {
			return c.Send("Error creating directory: " + err.Error())
		}
	}

	project, err := b.store.CreateProject(name, projectPath, githubURL)
	if err != nil {
		return c.Send("Error saving project: " + err.Error())
	}

	us.Step = state.StepInProject
	us.ActiveProject = &project.ID
	us.Mode = project.Mode

	text := fmt.Sprintf("Project *%s* created!\nMode: *%s*\n\nSend a task or use the buttons below.",
		project.Name, project.Mode)
	return c.Send(text, projectContextKeyboard(project.Mode), tele.ModeMarkdown)
}
