package mysql

import (
	"blog_demo/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "tingyusleep.cn"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExist
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

// encryptPassword 给用户密码加密（再存入数据库）
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 注册
func Login(user *models.User) (err error) {
	oPassword := user.Password //用户登录的密码
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrUserNotExist
	}
	if err != nil { //查询数据库错误
		return err
	}

	//判断密码是够正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrInvalidParam
	}
	return
}

// GetUserByID 根据用户id获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}
