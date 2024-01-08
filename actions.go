package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func showIndexPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/index.html")
}

func showLoginPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/login.html")
}

func showRegistrationPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/register.html")
}

func showForgotPasswordPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/forgot-password.html")
}

func showResetPasswordPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/reset-password.html")
}

func showConfirmEmailPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/confirm-email.html")
}

func showResendConfirmationEmailPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/resend-confirmation-email.html")
}

func showChangePasswordPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/change-password.html")
}

func showChangeEmailPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/change-email.html")
}

func show404Page(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "templates/404.html")
}

func login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeLoginRequest(&request)
	request.Validate()

	user, err := findUserByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.isConfirmed() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please confirm your email address"})
		return
	}

	sessionToken, err := createSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	environment := "development"
	if gin.Mode() == gin.ReleaseMode {
		environment = "production"
	}

	c.SetCookie("session_token_"+environment, sessionToken.Token, 3600, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}

func logout(c *gin.Context) {
	environment := "development"
	if gin.Mode() == gin.ReleaseMode {
		environment = "production"
	}

	c.SetCookie("session_token_"+environment, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeRegisterRequest(&request)
	request.Validate()

	if !isValidEmail(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	if !isValidPassword(request.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	if !isValidConfirmPassword(request.Password, request.ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
		return
	}

	_, err := createUser(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	environment := "development"
	if gin.Mode() == gin.ReleaseMode {
		environment = "production"
	}

	c.SetCookie("session_token_"+environment, "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered"})
}

func forgotPassword(c *gin.Context) {
	var request ForgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeForgotPasswordRequest(&request)
	request.Validate()

	if !isValidEmail(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	user, err := findUserByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := sendEmail("forgotPassword", user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully sent forgot password email"})
}

func resetPassword(c *gin.Context) {
	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeResetPasswordRequest(&request)
	request.Validate()

	if !isValidPassword(request.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	if !isValidConfirmPassword(request.Password, request.ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
		return
	}

	user, err := findUserByResetPasswordToken(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updatePassword(user, request.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully reset password"})
}

func confirmEmail(c *gin.Context) {
	var request ConfirmEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeConfirmEmailRequest(&request)
	request.Validate()

	user, err := findUserByConfirmEmailToken(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := user.setConfirmed(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully confirmed email"})
}

func resendConfirmationEmail(c *gin.Context) {
	var request ResendConfirmationEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeResendConfirmationEmailRequest(&request)
	request.Validate()

	if !isValidEmail(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	user, err := findUserByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := sendEmail("confirmationEmail", user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully sent confirmation email"})
}

func changePassword(c *gin.Context) {
	var request ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeChangePasswordRequest(&request)
	request.Validate()

	if !isValidPassword(request.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	if !isValidConfirmPassword(request.Password, request.ConfirmPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords don't match"})
		return
	}

	user, err := findUserBySession(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updatePassword(user, request.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed password"})
}

func changeEmail(c *gin.Context) {
	var request ChangeEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sanatizeChangeEmailRequest(&request)
	request.Validate()

	if !isValidEmail(request.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	user, err := findUserBySession(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updateEmail(user, request.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed email"})
}
