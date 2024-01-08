package golangauthsample

import (
	"fmt"
	"runtime"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Confirmed bool `gorm:"not null;default:false"`
}

type Session struct {
	gorm.Model
	UserID uint
	Token string `gorm:"unique_index;not null"`
}


func initializeRoutes(router *gin.Engine) {
	// Handle the index route
	router.GET("/", showIndexPage)

	// Group auth related routes together
	authRoutes := router.Group("/auth")
	{

		authRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		authRoutes.POST("/login", ensureNotLoggedIn(), performLogin)

		authRoutes.GET("/logout", ensureLoggedIn(), logout)
		authRoutes.POST("/logout", ensureLoggedIn(), logout)

		authRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)
		authRoutes.POST("/register", ensureNotLoggedIn(), register)

		authRoutes.GET("/forgot-password", ensureNotLoggedIn(), showForgotPasswordPage)
		authRoutes.POST("/forgot-password", ensureNotLoggedIn(), forgotPassword)

		authRoutes.GET("/reset-password", ensureNotLoggedIn(), showResetPasswordPage)
		authRoutes.POST("/reset-password", ensureNotLoggedIn(), resetPassword)

		authRoutes.GET("/confirm-email", ensureNotLoggedIn(), showConfirmEmailPage)
		authRoutes.POST("/confirm-email", ensureNotLoggedIn(), confirmEmail)

		authRoutes.GET("/resend-confirmation-email", ensureNotLoggedIn(), showResendConfirmationEmailPage)
		authRoutes.POST("/resend-confirmation-email", ensureNotLoggedIn(), resendConfirmationEmail)

		authRoutes.GET("/change-password", ensureLoggedIn(), showChangePasswordPage)
		authRoutes.POST("/change-password", ensureLoggedIn(), changePassword)

		authRoutes.GET("/change-email", ensureLoggedIn(), showChangeEmailPage)
		authRoutes.POST("/change-email", ensureLoggedIn(), changeEmail)


	}

}

func sessionExists(c *gin.Context) bool {
	_, err := c.Cookie("session_token")
	return err == nil
}

func getDBInstance() *gorm.DB {
	// use db file instance in memory 
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func findUserBySession(c *gin.Context) (User, error) {
	db := getDBInstance()
	sessionToken, err := c.Cookie("session_token")
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