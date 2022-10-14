package authentication

import "sync"

type user struct {
	name  string
	token string
}

type AuthenticationManager struct {
	currentUser   *user
	userChangeMux sync.RWMutex
}

func NewAuthenticationManager() *AuthenticationManager {
	return new(AuthenticationManager)
}

func (am *AuthenticationManager) ClearCurrentUser() {
	am.userChangeMux.Lock()
	defer am.userChangeMux.Unlock()
	am.currentUser = nil
}

func (am *AuthenticationManager) GetCurrentUser() string {
	am.userChangeMux.RLock()
	defer am.userChangeMux.RUnlock()
	if am.currentUser == nil {
		return ""
	}
	return am.currentUser.name
}

func (am *AuthenticationManager) GetCurrentUserCredential() string {
	am.userChangeMux.RLock()
	defer am.userChangeMux.RUnlock()
	if am.currentUser == nil {
		return ""
	}
	return am.currentUser.token
}

func (am *AuthenticationManager) SetCurrentUser(username, token string) {
	am.userChangeMux.Lock()
	defer am.userChangeMux.Unlock()
	newUser := &user{
		name:  username,
		token: token,
	}
	am.currentUser = newUser
}
