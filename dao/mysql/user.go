package mysql

import (
	"blog_demo/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

const secret = "tingyusleep.cn"

var (
	ErrorUserExist    = errors.New("用户已存在")
	ErrorUserNotExist = errors.New("用户不存在")
	ErrorInvalidParam = errors.New("用户名或密码错误")
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码加密加密
	user.Password = encryptPassword(user.Password)

	//执行sql语句入库
	sqlStr := `insert into user (user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password //用户登录的密码
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil { //查询数据库错误
		return err
	}

	//判断密码是够正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidParam
	}
	return
}
