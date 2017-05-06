package appContext

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

// AppContext custom echo.Context
type AppContext struct {
	echo.Context
	cookieStore *sessions.CookieStore
}

// NewAppContext returns a new AppContext
func NewAppContext(c echo.Context) *AppContext {
	return &AppContext{c, nil}
}

// SetCookieStore cookieStore setter
func (a *AppContext) SetCookieStore(cs *sessions.CookieStore) {
	a.cookieStore = cs
}

// GetCookieStore cookieStore getter
func (a *AppContext) GetCookieStore() *sessions.CookieStore {
	return a.cookieStore
}
