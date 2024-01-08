package golangauthsample

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path string
	Handler gin.HandlerFunc
	Method string
}

var protectedRoutes = []Route{
	{Path: "/auth/logout", Handler: logout, Method: "POST"},
	{Path: "/auth/change-password", Handler: showChangePasswordPage, Method: "GET"},
	{Path: "/auth/change-password", Handler: changePassword, Method: "POST"},
	{Path: "/auth/change-email", Handler: showChangeEmailPage, Method: "GET"},
	{Path: "/auth/change-email", Handler: changeEmail, Method: "POST"},
}

var unprotectedRoutes = []Route{
	{Path: "/", Handler: showIndexPage, Method: "GET"},
	{Path: "/auth/login", Handler: showLoginPage, Method: "GET"},
	{Path: "/auth/login", Handler: performLogin, Method: "POST"},
	{Path: "/auth/register", Handler: showRegistrationPage, Method: "GET"},
	{Path: "/auth/register", Handler: register, Method: "POST"},
	{Path: "/auth/forgot-password", Handler: showForgotPasswordPage, Method: "GET"},
	{Path: "/auth/forgot-password", Handler: forgotPassword, Method: "POST"},
	{Path: "/auth/reset-password", Handler: showResetPasswordPage, Method: "GET"},
	{Path: "/auth/reset-password", Handler: resetPassword, Method: "POST"},
	{Path: "/auth/confirm-email", Handler: showConfirmEmailPage, Method: "GET"},
	{Path: "/auth/confirm-email", Handler: confirmEmail, Method: "POST"},
	{Path: "/auth/resend-confirmation-email", Handler: showResendConfirmationEmailPage, Method: "GET"},
	{Path: "/auth/resend-confirmation-email", Handler: resendConfirmationEmail, Method: "POST"},
}

func (p Route) Register(router *gin.Engine) {
	router.Handle(p.Method, p.Path, p.Handler)
}

func initializeRoutes(router *gin.Engine) {
	for _, route := range protectedRoutes {
		route.Register(router)
	}
	for _, route := range unprotectedRoutes {
		route.Register(router)
	}
}