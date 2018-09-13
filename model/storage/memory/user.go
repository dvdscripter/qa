package memory

import (
	"time"

	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/crypto/argon2"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) Login(login string, pass string) error {
	user, err := db.FindUserByEmail(login)
	if err != nil {
		argon2.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		return errors.Errorf("user or pass invalid")
	}
	if err := argon2.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return errors.Errorf("user or pass invalid")
	}
	return nil
}

func (db *DB) CreateUser(u model.User) (model.User, error) {
	if _, err := db.FindUserByNick(u.Nick); err == nil {
		return model.User{}, errors.Errorf("Cannot create user")
	}
	if _, err := db.FindUserByEmail(u.Email); err == nil {
		return model.User{}, errors.Errorf("Cannot create user")
	}

	u.ID = len(db.users) + 1
	u.Since = time.Now()
	if errs := u.Valid(); errs != nil {
		return model.User{}, model.ErrInvalidUser
	}
	encPass, err := model.GenPass(u.Password)
	if err != nil {
		return model.User{}, err
	}
	u.Password = encPass

	db.users = append(db.users, u)
	u.Password = ""

	return u, nil
}

func (db *DB) FindAllUser() ([]model.User, error) {
	return model.OmitPass(db.users), nil
}

func (db *DB) UpdateUser(u model.User) (model.User, error) {
	if errs := u.Valid(); errs != nil {
		return model.User{}, model.ErrInvalidUser
	}
	for i, user := range db.users {
		if u.ID == user.ID {
			if u.Password != "" {
				newPass, err := model.GenPass(u.Password)
				if err != nil {
					return model.User{}, errors.Errorf("Cannot update password")
				}
				db.users[i].Password = newPass
			}
			// db.users[i].Email = html.EscapeString(u.Email)
			db.users[i].Nick = u.Nick
			db.users[i].Avatar = u.Avatar
			u.Password = ""
			return u, nil
		}
	}
	return model.User{}, storage.ErrUserNotFound
}

func (db *DB) DeleteUser(id int) error {
	index := -1
	for i, user := range db.users {
		if id == user.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return storage.ErrUserNotFound
	}

	db.users = append(db.users[:index], db.users[index+1:]...)

	return nil
}

func (db *DB) FindUser(id int) (model.User, error) {
	for _, user := range db.users {
		if id == user.ID {
			user.Password = ""
			return user, nil
		}
	}
	return model.User{}, storage.ErrUserNotFound
}

func (db *DB) FindUserByEmail(email string) (model.User, error) {
	for _, user := range db.users {
		if email == user.Email {
			return user, nil
		}
	}
	return model.User{}, storage.ErrUserNotFound
}

func (db *DB) FindUserByNick(nick string) (model.User, error) {
	for _, user := range db.users {
		if nick == user.Nick {
			return user, nil
		}
	}
	return model.User{}, storage.ErrUserNotFound
}
