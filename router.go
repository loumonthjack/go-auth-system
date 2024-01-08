package main

import (
	"context"
	"log"
	"os"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)



func initializeRoutes(router *gin.Engine) {
	// Handle the index route
	router.GET("/", showIndexPage)
	ctx := context.Background()
	
    oidcConfig := OIDCConfig{
        ClientID:     os.Getenv("OIDC_CLIENT_ID"),
        ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("SERVER_URL") + "/callback",
        ProviderURL:  os.Getenv("OIDC_PROVIDER_URL"),
        Scopes:       []string{"profile", "email"},
    }

    provider, oauth2Config, err := initOIDCProvider(ctx, oidcConfig)
    if err != nil {
        log.Fatalf("Failed to initialize OIDC provider: %v", err)
    }

    verifier := provider.Verifier(&oidc.Config{ClientID: oidcConfig.ClientID})

	// Group auth related routes together
	authRoutes := router.Group("/auth")
	{

		authRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		authRoutes.POST("/login", ensureNotLoggedIn(), login)
		authRoutes.POST("/logout", ensureLoggedIn(), logout)

		authRoutes.GET("/oidc", startOIDCFlow(oauth2Config))
		authRoutes.GET("/callback", oidcCallback(oauth2Config, verifier))
		authRoutes.GET("/oidc-logout", oidcLogout(oauth2Config))

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
