package model

type RegisterReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Hardware string `json:"hardware" form:"hardware"`
}
