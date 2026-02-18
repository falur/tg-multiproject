package bot

import (
	"context"
	"sync"
	"time"

	"tg-multiproject/internal/config"
	"tg-multiproject/internal/state"
	"tg-multiproject/internal/storage"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	tele         *tele.Bot
	cfg          *config.Config
	store        *storage.Storage
	state        *state.Manager
	mu           sync.Mutex
	cancelClaude context.CancelFunc
}

func New(cfg *config.Config, store *storage.Storage, sm *state.Manager) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	teleBot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		tele:  teleBot,
		cfg:   cfg,
		store: store,
		state: sm,
	}

	teleBot.Use(authMiddleware(cfg.AllowedUserID))

	// Commands
	teleBot.Handle("/start", b.handleStart)
	teleBot.Handle("/cancel", b.handleCancel)

	// Callbacks
	teleBot.Handle(&btnMyProjects, b.handleMyProjects)
	teleBot.Handle(&btnCreateProject, b.handleCreateProject)
	teleBot.Handle(&btnModePlan, b.handleModePlan)
	teleBot.Handle(&btnModeEdit, b.handleModeEdit)
	teleBot.Handle(&btnSessions, b.handleSessions)
	teleBot.Handle(&btnCancel, b.handleCancelCallback)
	teleBot.Handle(&btnSkip, b.handleSkip)
	teleBot.Handle(&btnBack, b.handleBack)

	teleBot.Handle(&tele.InlineButton{Unique: "select_project"}, b.handleSelectProject)
	teleBot.Handle(&tele.InlineButton{Unique: "resume_session"}, b.handleResumeSession)

	// Text messages
	teleBot.Handle(tele.OnText, b.handleText)

	return b, nil
}

func (b *Bot) Start() {
	b.tele.Start()
}

func (b *Bot) Stop() {
	b.mu.Lock()
	if b.cancelClaude != nil {
		b.cancelClaude()
	}
	b.mu.Unlock()
	b.tele.Stop()
}
