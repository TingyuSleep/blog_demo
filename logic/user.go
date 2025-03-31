package logic

import (
	"blog_demo/dao/mysql"
	"blog_demo/models"
	"blog_demo/pkg/jwt"
	"blog_demo/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询出错
		return err
	}

	//2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//3.保存进数据库
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，所以能够拿到user.UserID
	if err = mysql.Login(user); err != nil {
		//登录失败
		return "", err
	}
	//生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
