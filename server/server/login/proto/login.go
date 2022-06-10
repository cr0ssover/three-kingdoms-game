package proto

type LoginRsp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
	UID      int    `json:"uid"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hardware string `json:"hardware"`
}
