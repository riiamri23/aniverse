package service

import "sync"

type TokenManager struct {
	tokens map[string]string
	mu     sync.RWMutex
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokens: make(map[string]string),
	}
}

func (tm *TokenManager) SetToken(provider, token string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tokens[provider] = token
}

func (tm *TokenManager) GetToken(provider string) string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.tokens[provider]
}
