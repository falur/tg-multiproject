package bot

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tg-multiproject/internal/claude"
	"tg-multiproject/internal/state"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleTaskSubmission(c tele.Context, us *state.UserState) error {
	if us.ActiveProject == nil {
		return c.Send("No project selected. Use /start.")
	}

	project, err := b.store.GetProject(*us.ActiveProject)
	if err != nil {
		return c.Send("Error loading project: " + err.Error())
	}

	statusMsg, err := b.tele.Send(c.Recipient(), "Starting Claude...", runningKeyboard())
	if err != nil {
		return err
	}

	us.Step = state.StepRunning
	us.LastMessageID = statusMsg.ID
	us.LastChatID = c.Chat().ID

	ctx, cancel := context.WithCancel(context.Background())

	b.mu.Lock()
	b.cancelClaude = cancel
	b.mu.Unlock()

	cfg := claude.RunConfig{
		Prompt:    c.Text(),
		CWD:       project.Path,
		Mode:      us.Mode,
		Binary:    b.cfg.ClaudeBinary,
		SessionID: us.SessionID,
	}

	events, errc := claude.Run(ctx, cfg)

	go b.processClaudeStream(ctx, cancel, c.Sender().ID, events, errc, statusMsg)

	return nil
}

func (b *Bot) processClaudeStream(
	ctx context.Context,
	cancel context.CancelFunc,
	userID int64,
	events <-chan claude.StreamEvent,
	errc <-chan error,
	statusMsg *tele.Message,
) {
	defer cancel()

	stored := tele.StoredMessage{
		MessageID: strconv.Itoa(statusMsg.ID),
		ChatID:    statusMsg.Chat.ID,
	}

	var buf strings.Builder
	buf.WriteString("Claude is working...\n\n")
	dirty := false
	var lastSessionID string
	var resultText string

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	flush := func() {
		if !dirty {
			return
		}
		text := buf.String()
		if len(text) > 4000 {
			text = "..." + text[len(text)-3900:]
		}
		_, _ = b.tele.Edit(&stored, text, runningKeyboard())
		dirty = false
	}

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				// Channel closed — check for errors
				flush()
				select {
				case err := <-errc:
					if err != nil {
						buf.WriteString("\n\nError: " + err.Error())
					}
				default:
				}
				b.finishStream(userID, &stored, &buf, lastSessionID, resultText)
				return
			}
			b.handleEvent(&ev, &buf, &lastSessionID, &resultText)
			dirty = true

		case <-ticker.C:
			flush()

		case <-ctx.Done():
			_, _ = b.tele.Edit(&stored, buf.String()+"\n\nCancelled.")
			b.resetAfterStream(userID)
			return
		}
	}
}

func (b *Bot) handleEvent(ev *claude.StreamEvent, buf *strings.Builder, sessionID *string, resultText *string) {
	if ev.SessionID != "" {
		*sessionID = ev.SessionID
	}

	switch ev.Type {
	case "assistant":
		if ev.Message != nil {
			for _, block := range ev.Message.Content {
				switch block.Type {
				case "text":
					buf.WriteString(block.Text)
				case "tool_use":
					fmt.Fprintf(buf, "\n[Tool: %s]\n", block.Name)
				}
			}
		}

	case "result":
		if ev.Result != nil {
			if ev.Result.SessionID != "" {
				*sessionID = ev.Result.SessionID
			}
			if ev.Result.Result != "" {
				*resultText = ev.Result.Result
			}
			if ev.Result.TotalCost > 0 {
				fmt.Fprintf(buf, "\n\nCost: $%.4f | Turns: %d", ev.Result.TotalCost, ev.Result.NumTurns)
			}
		}
	}
}

func (b *Bot) finishStream(userID int64, stored *tele.StoredMessage, buf *strings.Builder, sessionID, resultText string) {
	us := b.state.Get(userID)

	// Save session if we got one
	if sessionID != "" && us.ActiveProject != nil {
		summary := resultText
		if len(summary) > 200 {
			summary = summary[:200]
		}
		_ = b.store.SaveSession(*us.ActiveProject, sessionID, summary)
		us.SessionID = sessionID
	}

	// Build final text
	finalText := buf.String()
	if resultText != "" {
		finalText = resultText
		if sessionID != "" {
			finalText += fmt.Sprintf("\n\nSession: `%s`", sessionID)
		}
	}

	// If text is too long, send as a file
	if len(finalText) > 4096 {
		doc := &tele.Document{
			File:     tele.FromReader(bytes.NewReader([]byte(finalText))),
			FileName: "result.md",
		}
		chatID := stored.ChatID
		chat := &tele.Chat{ID: chatID}
		_, _ = b.tele.Send(chat, doc)
		_, _ = b.tele.Edit(stored, "Result sent as file.", projectContextKeyboard(us.Mode))
	} else {
		if len(finalText) > 4000 {
			finalText = finalText[:4000] + "..."
		}
		_, _ = b.tele.Edit(stored, finalText, projectContextKeyboard(us.Mode))
	}

	b.resetAfterStream(userID)
}

func (b *Bot) resetAfterStream(userID int64) {
	us := b.state.Get(userID)
	us.Step = state.StepInProject

	b.mu.Lock()
	b.cancelClaude = nil
	b.mu.Unlock()
}
