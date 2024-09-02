package telegram

import (
	"gopkg.in/telebot.v3"
	"sync"
)

type UserContextManager interface {
	SetTrace(ctx telebot.Context, step string) UserContextManager
	ResetUserContext(ctx telebot.Context) UserContextManager
	DeleteUserContext(ctx telebot.Context) UserContextManager
}

// 处理用户上下文
type userContext struct {
	users map[int64][]string
	mu    sync.RWMutex
}

// NewUserContextManager 创建用户上下文管理器的新实例
func NewUserContextManager() UserContextManager {
	return &userContext{
		users: make(map[int64][]string),
	}
}

// ResetUserContext 将用户的上下文重置为初始状态
func (m *userContext) ResetUserContext(ctx telebot.Context) UserContextManager {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[ctx.Sender().ID] = make([]string, 0)
	return m
}

func (m *userContext) SetTrace(ctx telebot.Context, step string) UserContextManager {
	userId := ctx.Sender().ID
	m.mu.RLock()
	if _, ok := m.users[userId]; !ok {
		m.users[userId] = make([]string, 0)
	}
	m.users[userId] = append(m.users[userId], step)
	m.mu.RUnlock()
	return m
}

func (m *userContext) GetTraces(ctx telebot.Context) []string {
	m.mu.RLock()
	if traces, ok := m.users[ctx.Sender().ID]; ok {
		m.mu.RUnlock()
		return traces
	}
	return nil
}

// DeleteUserContext 删除用户的上下文
func (m *userContext) DeleteUserContext(ctx telebot.Context) UserContextManager {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, ctx.Sender().ID)
	return m
}
