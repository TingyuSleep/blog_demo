package models

//定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//User ID 从请求中获取当前的用户
	PostID    int64 `json:"post_id,string" binding:"required"`                //帖子id
	Direction int8  `json:"direction,string" binding:"required,oneof=-1 0 1"` //赞成票1，取消投票0，反对票-1
	// oneof 是 validator库中的一个规则，参数必须是其中之一
}
