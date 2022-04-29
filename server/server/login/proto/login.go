package proto

type LoginRsp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `josn:"session"`
	UID      int    `json:"uid"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IP       string `josn:"ip"`
	Hardware int    `json:"hardware"`
}
