package main

import (
	"github.com/gin-gonic/gin"
)

func sessionExists(c *gin.Context) bool {
	_, err := c.Cookie("session_token")
	return err == nil
}

func getSessionToken(c *gin.Context) (string, error) {
	environment := getCurrentEnvironment()
	return c.Cookie("session_token_" + environment)
}

func getCurrentEnvironment() string {
	if gin.Mode() == gin.ReleaseMode {
		return "production"
	}
	return "development"
}

func findUserBySession(c *gin.Context) (User, error) {
	db := getDBInstance()

	sessionToken, err := getSessionToken(c)
	if err != nil {
		return User{}, err
	}

	user, err := getUserFromSession(db, sessionToken)
	if err != nil {
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
