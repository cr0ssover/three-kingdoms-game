package controller

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/net"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login/model"
	"github.com/cr0ssover/three-kingdoms-game/server/server/login/proto"
	"github.com/cr0ssover/three-kingdoms-game/server/server/models"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
	"github.com/mitchellh/mapstructure"
)

var DefaultAccount = &Account{}

type Account struct{}

func (a *Account) Router(r *net.Router) {
	g := r.NewGroup("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	loginReq := &proto.LoginReq{}
	// map转结构体
	mapstructure.Decode(req.Body.Msg, loginReq)
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username = ?", loginReq.Username).Get(user)
	if err != nil {
		logger.Warn("user表查询错误")
		rsp.Body.Code = errcode.DBError
		return
	}
	if !ok {
		rsp.Body.Code = errcode.UserNotExist
		return
	}
	// 验证密码
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		rsp.Body.Code = errcode.PwdIncorrect
		return
	}

	// 生成token
	token, _ := utils.Award(user.UId)
	rsp.Body.Code = errcode.OK
	loginRsp := &proto.LoginRsp{
		UID:      user.UId,
		Username: user.Username,
		Session:  token,
		Password: "",
	}
	rsp.Body.Msg = loginRsp

	// 记录历史登录信息
	loginHistory := &model.LoginHistory{
		UId:      user.UId,
		State:    model.Login,
		Ip:       req.Conn.Addr(),
		Hardware: loginReq.Hardware,
		CTime:    time.Now(),
	}
	_, err = db.Engine.Table(loginHistory).Insert(loginHistory)
	if err != nil {
		logger.Warn("login_historyt表数据插入数据失败，err: ", err)
		return
	}

	// 记录上一次登录信息
	loginLast := &model.LoginLast{}
	ok, err = db.Engine.Table(loginLast).Where("uid = ?", user.UId).Get(loginLast)
	if err != nil {
		logger.Warn("login_last表数据查询失败，err: ", err)
		return
	}

	loginLast.UId = user.UId
	loginLast.Ip = req.Conn.Addr()
	loginLast.Session = token
	loginLast.Hardware = loginReq.Hardware
	loginLast.LoginTime = time.Now()
	loginLast.IsLogout = model.Login
	if ok {
		_, err = db.Engine.Table(loginLast).Where("uid = ?", user.UId).Update(loginLast)
		if err != nil {
			logger.Warn("login_last表数据更新数据失败,err: ", err)
		}
	} else {
		_, err := db.Engine.Table(loginLast).Insert(loginLast)
		if err != nil {
			logger.Warn("login_last表数据插入数据失败,err: ", err)
		}
	}

	// 缓存
	net.Mgr.UserLogin(req.Conn, user.UId, token)
}
