package controllers

import (
	"context"
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

	if err := gothic.StoreInSession("username", u.Email, c.Request, c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
