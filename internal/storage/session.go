package storage

import "time"

type Session struct {
	ID        int64     `db:"id"`
	ProjectID int64     `db:"project_id"`
	SessionID string    `db:"session_id"`
	Summary   string    `db:"summary"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Storage) SaveSession(projectID int64, sessionID, summary string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (project_id, session_id, summary) VALUES (?, ?, ?)",
		projectID, sessionID, summary,
	)
	return err
}

func (s *Storage) ListSessions(projectID int64) ([]Session, error) {
	var sessions []Session
	err := s.db.Select(&sessions,
		"SELECT * FROM sessions WHERE project_id = ? ORDER BY created_at DESC LIMIT 20",
		projectID,
	)
	return sessions, err
}

func (s *Storage) GetSession(id int64) (*Session, error) {
	var sess Session
	err := s.db.Get(&sess, "SELECT * FROM sessions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}
