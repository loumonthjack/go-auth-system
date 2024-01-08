package golangauthsample

import (
	"github.com/gin-gonic/gin"
)

func sessionExists(c *gin.Context) bool {
	_, err := c.Cookie("session_token")
	return err == nil
}

func findUserBySession(c *gin.Context) (User, error) {
	db := getDBInstance()
	environment := "development"
	if gin.Mode() == gin.ReleaseMode {
		environment = "production"
	}
	sessionToken, err := c.Cookie("session_token_" + environment)
	if err != nil {
		return User{}, err
	}

	var session Session
	if err := db.Where("token = ?", sessionToken).First(&session).Error; err != nil {
		return User{}, err
	}

	var user User
	if err := db.Where("id = ?", session.UserID).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessionExists(c) {
			c.Redirect(302, "/")
			c.Abort()
			return
		}
	}
}

func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If there's no session cookie, redirect the user to the login page
		if !sessionExists(c) {
			c.Redirect(302, "/auth/login")
			c.Abort()
			return
		}

		// If the user exists in the database, proceed as normal
		if _, err := findUserBySession(c); err == nil {
			return
		}

		// If the user doesn't exist in the database, redirect the user to the login page
		c.Redirect(302, "/auth/login")
		c.Abort()
	}
}