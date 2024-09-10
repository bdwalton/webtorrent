package controllers

import (
	"context"
	"fmt"
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

func SignIn(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sign in with Google</title>
</head>
<body
     style="
         font-family: Arial, sans-serif;
         background-color: #f0f0f0;
         display: flex;
         justify-content: center;
         align-items: center;
         height: 100vh;"
>
  <div
     style="
        background-color: #fff;
        padding: 40px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        text-align: center;"
    >
        <h1 style="color: #333; margin-bottom: 20px;">WebTorrent Signin</h1>
        <a href="auth/google"  style="
        display: inline-flex;
        align-items: center;
        justify-content: center;
        background-color: #4285f4;
        color: #fff;
        text-decoration: none;
        padding: 12px 20px;
        border-radius: 4px;
        transition: background-color 0.3s ease;">
        <span style="font-size: 16px; font-weight: bold;">Sign in with Google</span>
        </a>
  </div>
</body>
</html>
`)))
}
