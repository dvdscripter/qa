package sql

import (
	"fmt"
	"html"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) FindAllUser() ([]model.User, error) {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return model.OmitPass(users), nil
}

func (db *DB) CreateUser(u model.User) (model.User, error) {
	if _, err := db.FindUserByNick(u.Nick); err == nil {
		return model.User{}, fmt.Errorf("%s already exist", u.Nick)
	}
	if _, err := db.FindUserByEmail(u.Email); err == nil {
		return model.User{}, fmt.Errorf("%s already exist", u.Email)
	}

	u.Since = time.Now()
	u.Email = html.EscapeString(u.Email)
	if errs := u.Valid(); errs != nil {
		return model.User{}, model.ErrInvalidUser
	}
	encPass, err := model.GenPass(u.Password)
	if err != nil {
		return model.User{}, err
	}
	u.Password = encPass

	if err := db.Create(&u).Error; err != nil {
		return model.User{}, err
	}
	u.Password = ""

	return u, nil

}

func (db *DB) UpdateUser(u model.User) (model.User, error) {
	if errs := u.Valid(); errs != nil {
		return model.User{}, model.ErrInvalidUser
	}

	user, err := db.FindUser(u.ID)
	if err != nil {
		return model.User{}, storage.ErrUserNotFound
	}

	if u.Password != "" {
		newPass, err := model.GenPass(u.Password)
		if err != nil {
			return model.User{}, fmt.Errorf("Cannot update password")
		}
		user.Password = newPass
	}
	// user.Email = html.EscapeString(u.Email)
	user.Nick = u.Nick
	user.Avatar = html.EscapeString(u.Avatar)
	if err := db.Save(&user).Error; err != nil {
		return model.User{}, err
	}
	u.Password = ""
	u.Since = user.Since

	return u, nil
}

func (db *DB) DeleteUser(id int) error {
	return db.Where("ID = ?", id).Delete(&model.User{}).Error
}

func (db *DB) FindUser(id int) (model.User, error) {
	var user model.User

	if err := db.First(&user, id).Error; err != nil {
		return model.User{}, storage.ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

func (db *DB) FindUserByNick(nick string) (model.User, error) {
	var user model.User

	if err := db.Where("nick = ?", nick).First(&user).Error; err != nil {
		return model.User{}, storage.ErrUserNotFound
	}

	return user, nil
}

func (db *DB) FindUserByEmail(email string) (model.User, error) {
	var user model.User

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, storage.ErrUserNotFound
	}

	return user, nil
}

func (db *DB) Login(login string, pass string) error {
	user, err := db.FindUserByEmail(login)
	if err != nil {
		return fmt.Errorf("user or pass invalid")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return fmt.Errorf("user or pass invalid")
	}
	return nil
}
