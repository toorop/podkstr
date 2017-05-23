package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"net/url"

	"encoding/base64"
	"encoding/hex"

	"github.com/labstack/echo"
	"github.com/toorop/podkstr/core"
	"github.com/toorop/podkstr/logger"
)

// DiscourseSSO discourse board SSO
// https://meta.discourse.org/t/official-single-sign-on-for-discourse/13045
func DiscourseSSO(c echo.Context) error {
	const URL = "https://board.podkstr.com/session/sso_login?"
	const secret = "jeuoapnsbgettdjsssqhdk"
	var nonce string
	// chek auth
	u := c.Get("user")
	if u == nil {
		return c.String(http.StatusOK, "Please first login into podkstr at https://podkstr.com/signin and retry board login")
	}
	// Parse query string
	payload64 := c.QueryParam("sso")
	sig := c.QueryParam("sig")

	// chek sig
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(payload64))
	if err != nil {
		logger.Log.Errorf("%s - discourseSSO -> unable to hash sso %s - %s ", c.Request().RemoteAddr, payload64, err)
		return c.NoContent(http.StatusInternalServerError)
	}
	sig2 := hex.EncodeToString(h.Sum(nil))
	if sig != sig2 {
		logger.Log.Infof("%s - discourseSSO -> bad sig ", c.Request().RemoteAddr)
		return c.NoContent(http.StatusForbidden)
	}

	// Get nonce
	payload, err := base64.StdEncoding.DecodeString(payload64)
	if err != nil {
		logger.Log.Errorf("%s - discourseSSO -> unable to base64 decode payload %s - %s ", c.Request().RemoteAddr, payload64, err)
		return c.NoContent(http.StatusInternalServerError)
	}
	queryStr, err := url.ParseQuery(string(payload))
	if err != nil {
		logger.Log.Errorf("%s - discourseSSO -> unable to parse payload %s - %s ", c.Request().RemoteAddr, payload, err)
		return c.NoContent(http.StatusInternalServerError)
	}
	nonce = queryStr.Get("nonce")

	// login
	// assume user is already logged for a first try
	user := u.(core.User)
	q := make(url.Values)
	q.Set("nonce", nonce)
	q.Set("email", user.Email)
	q.Set("external_id", user.UUID)

	q64 := base64.StdEncoding.EncodeToString([]byte(q.Encode()))

	h.Reset()
	_, err = h.Write([]byte(q64))
	if err != nil {
		logger.Log.Errorf("%s - discourseSSO -> unable to hash new sso %s - %s ", c.Request().RemoteAddr, q64, err)
		return c.NoContent(http.StatusInternalServerError)
	}
	sig = hex.EncodeToString(h.Sum(nil))

	q2 := make(url.Values)
	q2.Set("sso", q64)
	q2.Set("sig", sig)

	redirectURL := URL + q2.Encode()
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
