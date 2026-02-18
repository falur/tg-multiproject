package bot

import (
	"fmt"

	"tg-multiproject/internal/storage"

	tele "gopkg.in/telebot.v3"
)

var (
	btnMyProjects    = tele.InlineButton{Unique: "my_projects", Text: "My Projects"}
	btnCreateProject = tele.InlineButton{Unique: "create_project", Text: "Create Project"}
	btnModePlan      = tele.InlineButton{Unique: "mode_plan", Text: "Plan Mode"}
	btnModeEdit      = tele.InlineButton{Unique: "mode_edit", Text: "Edit Mode"}
	btnSessions      = tele.InlineButton{Unique: "sessions", Text: "Sessions"}
	btnCancel        = tele.InlineButton{Unique: "cancel_task", Text: "Cancel"}
	btnSkip          = tele.InlineButton{Unique: "skip", Text: "Skip"}
	btnBack          = tele.InlineButton{Unique: "back", Text: "Back"}
)

func mainMenuKeyboard() *tele.ReplyMarkup {
	rm := &tele.ReplyMarkup{}
	rm.InlineKeyboard = [][]tele.InlineButton{
		{btnMyProjects, btnCreateProject},
	}
	return rm
}

func projectListKeyboard(projects []storage.Project) *tele.ReplyMarkup {
	rm := &tele.ReplyMarkup{}
	for _, p := range projects {
		btn := tele.InlineButton{
			Unique: "select_project",
			Text:   p.Name,
			Data:   fmt.Sprintf("%d", p.ID),
		}
		rm.InlineKeyboard = append(rm.InlineKeyboard, []tele.InlineButton{btn})
	}
	rm.InlineKeyboard = append(rm.InlineKeyboard, []tele.InlineButton{btnBack})
	return rm
}

func projectContextKeyboard(mode string) *tele.ReplyMarkup {
	rm := &tele.ReplyMarkup{}
	var modeBtn tele.InlineButton
	if mode == "plan" {
		modeBtn = btnModeEdit
		modeBtn.Text = "Switch to Edit Mode"
	} else {
		modeBtn = btnModePlan
		modeBtn.Text = "Switch to Plan Mode"
	}
	rm.InlineKeyboard = [][]tele.InlineButton{
		{modeBtn},
		{btnSessions},
		{btnBack},
	}
	return rm
}

func runningKeyboard() *tele.ReplyMarkup {
	rm := &tele.ReplyMarkup{}
	rm.InlineKeyboard = [][]tele.InlineButton{
		{btnCancel},
	}
	return rm
}

func sessionListKeyboard(sessions []storage.Session) *tele.ReplyMarkup {
	rm := &tele.ReplyMarkup{}
	for _, s := range sessions {
		summary := s.Summary
		if summary == "" {
			summary = s.SessionID[:min(len(s.SessionID), 16)]
		}
		if len(summary) > 40 {
			summary = summary[:40] + "..."
		}
		btn := tele.InlineButton{
			Unique: "resume_session",
			Text:   summary,
			Data:   fmt.Sprintf("%d", s.ID),
		}
		rm.InlineKeyboard = append(rm.InlineKeyboard, []tele.InlineButton{btn})
	}
	rm.InlineKeyboard = append(rm.InlineKeyboard, []tele.InlineButton{btnBack})
	return rm
}
