package controllers

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func SignInWithProvider(c *gin.Context) {
	r := c.Request.Clone(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
	gothic.BeginAuthHandler(c.Writer, r)
}

func CallBackHandler(c *gin.Context) {
	r := c.Request.Clone(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
	u, err := gothic.CompleteUserAuth(c.Writer, r)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, _ := store.Get(c.Request, "webtorrent-session")
	session.Values["username"] = u.Email
	session.Save(r, c.Writer)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
