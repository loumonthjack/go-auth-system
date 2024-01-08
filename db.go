package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique_index;not null"`
	Password  string `gorm:"not null"`
	Confirmed bool   `gorm:"not null;default:false"`
}

type Session struct {
	gorm.Model
	UserID uint
	Token  string `gorm:"unique_index;not null"`
}

func (u *User) isConfirmed() bool {
	return u.Confirmed
}

func (u *User) setConfirmed() error {
	u.Confirmed = true
	db := getDBInstance()
	if err := db.Save(&u).Error; err != nil {
		return err
	}
	return nil
}

// initiateDB creates a new connection to our postgres database.

func initiateDB() *gorm.DB {
	dsn := "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func getDBInstance() *gorm.DB {
	channel := make(chan *gorm.DB, 1)
	if channel != nil {
		return <-channel
	}
	go func() {
		channel <- initiateDB()
		models := []interface{}{
			&User{},
			&Session{},
		}
		db := <-channel
		db.AutoMigrate(models...)
		db.Debug()
		channel <- db
	}()

	return <-channel

}

func createSession(userID uint) (Session, error) {
	db := getDBInstance()
	session := Session{UserID: userID}
	if err := db.Create(&session).Error; err != nil {
		return Session{}, err
	}
	return session, nil
}

func findUserByEmail(email string) (User, error) {
	db := getDBInstance()
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func createUser(email string, password string) (User, error) {
	db := getDBInstance()
	user := User{Email: email, Password: password}
	if err := db.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func updatePassword(user User, password string) error {
	db := getDBInstance()
	if err := db.Model(&user).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func findUserByConfirmEmailToken(token string) (User, error) {
	db := getDBInstance()
	var user User
	if err := db.Where("confirm_email_token = ?", token).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func findUserByResetPasswordToken(token string) (User, error) {
	db := getDBInstance()
	var user User
	if err := db.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func findSessionByToken(token string) (Session, error) {
	db := getDBInstance()
	var session Session
	if err := db.Where("token = ?", token).First(&session).Error; err != nil {
		return Session{}, err
	}
	return session, nil
}

func deleteSession(session Session) error {
	db := getDBInstance()
	if err := db.Delete(&session).Error; err != nil {
		return err
	}
	return nil
}

func updateEmail(user User, email string) error {
	db := getDBInstance()
	if err := db.Model(&user).Update("email", email).Error; err != nil {
		return err
	}
	return nil
}

func getUserFromSession(db *gorm.DB, sessionToken string) (User, error) {
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
