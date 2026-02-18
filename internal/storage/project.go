package storage

import "time"

type Project struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	GitHubURL string    `db:"github_url"`
	Mode      string    `db:"mode"`
	CreatedAt time.Time `db:"created_at"`
}

func (s *Storage) CreateProject(name, path, githubURL string) (*Project, error) {
	res, err := s.db.Exec(
		"INSERT INTO projects (name, path, github_url) VALUES (?, ?, ?)",
		name, path, githubURL,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return s.GetProject(id)
}

func (s *Storage) GetProject(id int64) (*Project, error) {
	var p Project
	err := s.db.Get(&p, "SELECT * FROM projects WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Storage) ListProjects() ([]Project, error) {
	var projects []Project
	err := s.db.Select(&projects, "SELECT * FROM projects ORDER BY created_at DESC")
	return projects, err
}

func (s *Storage) UpdateProjectMode(id int64, mode string) error {
	_, err := s.db.Exec("UPDATE projects SET mode = ? WHERE id = ?", mode, id)
	return err
}
