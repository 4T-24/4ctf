package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/gorilla/securecookie"
	"github.com/valyala/fasthttp"
)

// SessionManager manages type-safe session data using securecookie.
type SessionManager[Session any] struct {
	secureCookie *securecookie.SecureCookie
	sessionName  string
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager[Session any](hashKey, blockKey []byte, sessionName string, sessionType Session) *SessionManager[Session] {
	secureCookie := securecookie.New(hashKey, blockKey)
	return &SessionManager[Session]{
		secureCookie: secureCookie,
		sessionName:  sessionName,
	}
}

// encodeSession serializes session data using gob.
func (sm *SessionManager[Session]) encodeSession(data *Session) (string, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// decodeSession deserializes session data using gob.
func (sm *SessionManager[Session]) decodeSession(encoded string) (*Session, error) {
	var data Session
	decoder := gob.NewDecoder(bytes.NewBufferString(encoded))
	if err := decoder.Decode(&data); err != nil {
		return new(Session), err
	}
	return &data, nil
}

// SetSession sets the session data for a user.
func (sm *SessionManager[Session]) SetSession(ctx *fasthttp.RequestCtx, data *Session) error {
	encodedData, err := sm.encodeSession(data)
	if err != nil {
		return err
	}

	secureValue, err := sm.secureCookie.Encode(sm.sessionName, encodedData)
	if err != nil {
		return err
	}

	cookie := &fasthttp.Cookie{}
	cookie.SetKey(sm.sessionName)
	cookie.SetValue(secureValue)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true) // Adjust based on your environment
	ctx.Response.Header.SetCookie(cookie)
	return nil
}

// GetSession retrieves session data for a user.
func (sm *SessionManager[Session]) GetSession(ctx *fasthttp.RequestCtx) (*Session, error) {
	cookie := ctx.Request.Header.Cookie(sm.sessionName)
	if cookie == nil {
		return new(Session), fmt.Errorf("no session cookie found")
	}

	var encodedData string
	if err := sm.secureCookie.Decode(sm.sessionName, string(cookie), &encodedData); err != nil {
		return new(Session), err
	}

	data, err := sm.decodeSession(encodedData)
	if err != nil {
		return new(Session), err
	}
	return data, nil
}

// DeleteSession removes the session for a user.
func (sm *SessionManager[Session]) DeleteSession(ctx *fasthttp.RequestCtx) {
	cookie := &fasthttp.Cookie{}
	cookie.SetKey(sm.sessionName)
	cookie.SetValue("") // Clear the cookie value
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true) // Adjust based on your environment
	cookie.SetExpire(fasthttp.CookieExpireDelete)
	ctx.Response.Header.SetCookie(cookie)
}
