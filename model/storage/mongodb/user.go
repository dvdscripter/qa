package mongodb

import (
	"html"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/crypto/argon2"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) FindAllUser() ([]model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	var users []model.User
	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).Find(nil).All(&users); err != nil {
		return nil, errors.Wrap(err, "cannot enumerate users")
	}
	return model.OmitPass(users), nil
}

func (db *DB) CreateUser(u model.User) (model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	if _, err := db.FindUserByNick(u.Nick); err == nil {
		return model.User{}, errors.Errorf("Cannot create user")
	}
	if _, err := db.FindUserByEmail(u.Email); err == nil {
		return model.User{}, errors.Errorf("Cannot create user")
	}

	u.Since = time.Now()
	if errs := u.Valid(); errs != nil {
		return model.User{}, errors.Errorf("Cannot create user: %s", errs)
	}
	encPass, err := model.GenPass(u.Password)
	if err != nil {
		return model.User{}, err
	}
	u.Password = encPass

	u.ID = db.getID(db.GetUserC())

	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).Insert(&u); err != nil {
		return model.User{}, errors.Wrap(err, "cannot create new user")
	}
	u.Password = ""

	return u, nil
}

func (db *DB) UpdateUser(u model.User) (model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	user, err := db.FindUser(u.ID)
	if err != nil {
		return model.User{}, storage.ErrUserNotFound
	}

	if u.Password != "" {
		if errs := u.ValidPassword(); errs != nil {
			return model.User{}, errs
		}
		newPass, err := model.GenPass(u.Password)
		if err != nil {
			return model.User{}, err
		}
		user.Password = newPass
	}

	if errs := u.ValidNick(); errs != nil {
		return model.User{}, errs
	}
	if errs := u.ValidAvatar(); errs != nil {
		return model.User{}, errs
	}

	user.Nick = u.Nick
	user.Avatar = html.EscapeString(u.Avatar)

	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).UpdateId(user.ID, &user); err != nil {
		return model.User{}, errors.Wrapf(err, "cannot update user: %s", user.Nick)
	}
	u.Password = ""
	u.Since = user.Since

	return u, nil
}

func (db *DB) DeleteUser(id int) error {
	conn := db.Copy()
	defer conn.Close()

	return conn.DB(db.GetDatabase()).C(db.GetUserC()).RemoveId(id)
}

func (db *DB) FindUser(id int) (model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	var user model.User

	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).FindId(id).One(&user); err != nil {
		return model.User{}, storage.ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

func (db *DB) FindUserByNick(nick string) (model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	var user model.User

	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).Find(bson.M{"nick": nick}).One(&user); err != nil {
		return model.User{}, storage.ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

func (db *DB) FindUserByEmail(email string) (model.User, error) {
	conn := db.Copy()
	defer conn.Close()

	var user model.User

	if err := conn.DB(db.GetDatabase()).C(db.GetUserC()).Find(bson.M{"email": email}).One(&user); err != nil {
		return model.User{}, storage.ErrUserNotFound
	}

	return user, nil
}

func (db *DB) Login(login string, pass string) error {
	user, err := db.FindUserByEmail(login)
	if err != nil {
		argon2.CompareHashAndPassword([]byte("not found"), []byte(pass))
		return errors.Errorf("user or pass invalid")
	}
	if err := argon2.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return errors.Errorf("user or pass invalid")
	}
	return nil
}
