package state

import "sync"

type Step int

const (
	StepIdle Step = iota
	StepCreateName
	StepCreateURL
	StepInProject
	StepRunning
)

type UserState struct {
	Step          Step
	ActiveProject *int64
	ProjectName   string
	Mode          string
	SessionID     string
	LastMessageID int
	LastChatID    int64
}

type Manager struct {
	mu     sync.Mutex
	states map[int64]*UserState
}

func NewManager() *Manager {
	return &Manager{
		states: make(map[int64]*UserState),
	}
}

func (m *Manager) Get(userID int64) *UserState {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.states[userID]
	if !ok {
		s = &UserState{Step: StepIdle}
		m.states[userID] = s
	}
	return s
}

func (m *Manager) Set(userID int64, s *UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[userID] = s
}

func (m *Manager) Reset(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[userID] = &UserState{Step: StepIdle}
}
