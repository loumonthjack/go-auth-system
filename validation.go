package golangauthsample

import (
	"strings"
	"errors"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (u *User) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func isValidEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 12
}

func isValidConfirmPassword(password string, confirmPassword string) bool {
	return password == confirmPassword
}


type RegisterRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirmPassword"`
}

func sanatizeRegisterRequest(newUser *RegisterRequest) {
	newUser.Email = strings.ToLower(newUser.Email)
	newUser.Email = govalidator.Trim(newUser.Email, "")
	newUser.Password = govalidator.Trim(newUser.Password, "")
	newUser.ConfirmPassword = govalidator.Trim(newUser.ConfirmPassword, "")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func sanatizeLoginRequest(loginRequest *LoginRequest) {
	loginRequest.Email = strings.ToLower(loginRequest.Email)
	loginRequest.Email = govalidator.Trim(loginRequest.Email, "")
	loginRequest.Password = govalidator.Trim(loginRequest.Password, "")
}

func InvalidMessage(message string) error {
	return errors.New(message)
}

func (u *RegisterRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidEmail(u.Email) {
		return InvalidMessage("Invalid Email")
	}

	if !isValidPassword(u.Password) {
		return InvalidMessage("Invalid Password")
	}

	if !isValidConfirmPassword(u.Password, u.ConfirmPassword) {
		return InvalidMessage("Passwords do not match")
	}

	return nil
}
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func sanatizeForgotPasswordRequest(forgotPasswordRequest *ForgotPasswordRequest) {
	forgotPasswordRequest.Email = strings.ToLower(forgotPasswordRequest.Email)
	forgotPasswordRequest.Email = govalidator.Trim(forgotPasswordRequest.Email, "")
}

func (u *ForgotPasswordRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidEmail(u.Email) {
		return InvalidMessage("Invalid Email")
	}

	return nil
}

type ResetPasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Token           string `json:"token"`
}

func sanatizeResetPasswordRequest(resetPasswordRequest *ResetPasswordRequest) {
	resetPasswordRequest.Password = govalidator.Trim(resetPasswordRequest.Password, "")
	resetPasswordRequest.ConfirmPassword = govalidator.Trim(resetPasswordRequest.ConfirmPassword, "")
	resetPasswordRequest.Token = govalidator.Trim(resetPasswordRequest.Token, "")
}

func (u *ResetPasswordRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidPassword(u.Password) {
		return InvalidMessage("Invalid Password")
	}

	if !isValidConfirmPassword(u.Password, u.ConfirmPassword) {
		return InvalidMessage("Passwords do not match")
	}

	return nil
}

type UpdatePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func sanatizeUpdatePasswordRequest(updatePasswordRequest *UpdatePasswordRequest) {
	updatePasswordRequest.Password = govalidator.Trim(updatePasswordRequest.Password, "")
	updatePasswordRequest.ConfirmPassword = govalidator.Trim(updatePasswordRequest.ConfirmPassword, "")
}

func (u *UpdatePasswordRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidPassword(u.Password) {
		return InvalidMessage("Invalid Password")
	}

	if !isValidConfirmPassword(u.Password, u.ConfirmPassword) {
		return InvalidMessage("Passwords do not match")
	}

	return nil
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func sanatizeUpdateEmailRequest(updateEmailRequest *UpdateEmailRequest) {
	updateEmailRequest.Email = strings.ToLower(updateEmailRequest.Email)
	updateEmailRequest.Email = govalidator.Trim(updateEmailRequest.Email, "")
}

func (u *UpdateEmailRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidEmail(u.Email) {
		return InvalidMessage("Invalid Email")
	}

	return nil
}

type ChangePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func sanatizeChangePasswordRequest(changePasswordRequest *ChangePasswordRequest) {
	changePasswordRequest.Password = govalidator.Trim(changePasswordRequest.Password, "")
	changePasswordRequest.ConfirmPassword = govalidator.Trim(changePasswordRequest.ConfirmPassword, "")
}

func (u *ChangePasswordRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidPassword(u.Password) {
		return InvalidMessage("Invalid Password")
	}

	if !isValidConfirmPassword(u.Password, u.ConfirmPassword) {
		return InvalidMessage("Passwords do not match")
	}

	return nil
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

func sanatizeChangeEmailRequest(changeEmailRequest *ChangeEmailRequest) {
	changeEmailRequest.Email = strings.ToLower(changeEmailRequest.Email)
	changeEmailRequest.Email = govalidator.Trim(changeEmailRequest.Email, "")
}

func (u *ChangeEmailRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidEmail(u.Email) {
		return InvalidMessage("Invalid Email")
	}

	return nil
}

type ConfirmEmailRequest struct {
	Token string `json:"token"`
}

func sanatizeConfirmEmailRequest(confirmEmailRequest *ConfirmEmailRequest) {
	confirmEmailRequest.Token = govalidator.Trim(confirmEmailRequest.Token, "")
}

func (u *ConfirmEmailRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	return nil
}

type SessionRequest struct {
	Token string `json:"token"`
}

func sanatizeSessionRequest(sessionRequest *SessionRequest) {
	sessionRequest.Token = govalidator.Trim(sessionRequest.Token, "")
}

func (u *SessionRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	return nil
}

type ResendConfirmationEmailRequest struct {
	Email string `json:"email"`
}

func sanatizeResendConfirmationEmailRequest(resendConfirmationEmailRequest *ResendConfirmationEmailRequest) {
	resendConfirmationEmailRequest.Email = strings.ToLower(resendConfirmationEmailRequest.Email)
	resendConfirmationEmailRequest.Email = govalidator.Trim(resendConfirmationEmailRequest.Email, "")
}

func (u *ResendConfirmationEmailRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return err
	}

	if !isValidEmail(u.Email) {
		return InvalidMessage("Invalid Email")
	}

	return nil
}