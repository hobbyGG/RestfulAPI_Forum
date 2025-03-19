package mysql

import (
	"github.com/hobbyGG/RestfulAPI_forum/models"
	"go.uber.org/zap"
)

func AddUser(newUser *models.User) error {
	sqlStr := `
	insert into user(user_id, username, password, email)
	values(?, ?, ?, ?)`
	_, err := db.Exec(sqlStr,
		newUser.UID,
		newUser.UserName,
		newUser.Pwd,
		newUser.Email)
	if err != nil {
		zap.L().Error("AddUser exec error", zap.Error(err))
		return err
	}
	return nil
}

func GetUserByID(id int64) (*models.User, error) {
	sqlStr := `
	select user_id, username, password, email from user where user_id=?`
	var user models.User
	if err := db.Get(&user, sqlStr, id); err != nil {
		zap.L().Error("db.Get error", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

// 删除一个用户并放回该用户数据
func DelUserByID(id int64) (*models.User, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return nil, err
	}
	sqlStr := `
	delete from user where user_id=?`
	_, err = db.Exec(sqlStr, id)
	if err != nil {
		zap.L().Error("delete exec error", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func GetUsers(username string) ([]*models.User, error) {
	users := make([]*models.User, 1)
	sqlStr := `
	select user_id, username, password, email from user
	where username = ?`
	err := db.Select(&users, sqlStr, username)
	if err != nil {
		zap.L().Error("db.Select error", zap.Error(err))
		return nil, err
	}
	return users, nil
}
