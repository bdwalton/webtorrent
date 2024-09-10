package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

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

	if _, ok := allowedUsers[u.Email]; !ok {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("User %q not permitted.", u.Email))
		return
	}

	if err := gothic.StoreInSession("username", u.Email, c.Request, c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func SignoutHandler(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
