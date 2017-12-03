/*
	Digivance MVC Application Framework
	Session Manager Features
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for an in process browser session manager system. (E.g. per user
	server side memory map)
*/

package mvcapp

import (
	"time"
)

// SessionManager is the base struct that manages the collection
// of current http session models.
type SessionManager struct {
	// SessionIDKey is the name of the cookie value that will store the unique ID of the browser
	// session
	SessionIDKey string

	// Sessions is the collection of browser session objects
	Sessions []*Session

	// SessionTimeout is the duration of time that a browser session will stay in memory between
	// requests / activity from the user
	SessionTimeout time.Duration
}

// NewSessionManager returns a new Session Manager object
func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions:       make([]*Session, 0),
		SessionTimeout: (15 * time.Minute),
	}
}

// GetSession returns the current http session for the provided session id
func (manager *SessionManager) GetSession(id string) *Session {
	for key, val := range manager.Sessions {
		if val.ID == id {
			return manager.Sessions[key]
		}
	}

	return nil
}

// Contains detects if the requested id (key) exists in this session collection
func (manager *SessionManager) Contains(id string) bool {
	for _, v := range manager.Sessions {
		if v.ID == id {
			return true
		}
	}

	return false
}

// CreateSession creates and returns a new http session model
func (manager *SessionManager) CreateSession(id string) *Session {
	i := len(manager.Sessions)
	session := NewSession()
	session.ID = id
	manager.Sessions = append(manager.Sessions, session)
	return manager.Sessions[i]
}

// SetSession will set (creating if necessary) the provided session to
// the session manager collection
func (manager *SessionManager) SetSession(session *Session) {
	id := session.ID
	res := manager.GetSession(id)

	if res != nil {
		res.Values = append([]*SessionValue{}, session.Values...)
	} else {
		manager.Sessions = append(manager.Sessions, session)
	}
}

// DropSession will remove a session from the session manager collection based
// on the provided session id
func (manager *SessionManager) DropSession(id string) {
	for key, val := range manager.Sessions {
		if val.ID == id {
			if key > 1 {
				manager.Sessions = append(manager.Sessions[:key], manager.Sessions[key+1:]...)
			} else {
				if key == 1 {
					manager.Sessions = append(manager.Sessions[2:], manager.Sessions[0])
				} else {
					manager.Sessions = manager.Sessions[1:]
				}
			}
		}
	}
}

// CleanSessions will drop inactive sessions
func (manager *SessionManager) CleanSessions() {
	expired := time.Now().Add(-manager.SessionTimeout)

	for key, val := range manager.Sessions {
		if val.ActivityDate.Before(expired) {
			if key > 1 {
				if len(manager.Sessions) > 1 {
					manager.Sessions = append(manager.Sessions[:key], manager.Sessions[key+1:]...)
				} else {
					manager.Sessions = manager.Sessions[:key]
				}
			} else {
				if key == 1 {
					if len(manager.Sessions) > 1 {
						manager.Sessions = append(manager.Sessions[2:], manager.Sessions[0])
					} else {
						manager.Sessions = append([]*Session{}, manager.Sessions[0])
					}
				} else {
					if len(manager.Sessions) > 1 {
						manager.Sessions = manager.Sessions[1:]
					} else {
						manager.Sessions = []*Session{}
					}
				}
			}
		}
	}
}
