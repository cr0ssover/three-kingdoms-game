package logic

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/db"
	"github.com/cr0ssover/three-kingdoms-game/server/logger"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/models"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web/model"
	"github.com/cr0ssover/three-kingdoms-game/server/utils"
)

var DefaultAccountLogic = &AccountLogic{}

type AccountLogic struct {
}

func (*AccountLogic) Register(req *model.RegisterReq) error {
	username := req.Username
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		logger.Warn("注册数据查询失败,err: %v", err)
		return common.New(errcode.DBError, "数据库异常")
	}
	if ok {
		return common.New(errcode.UserExist, "用户已存在")
	}
	user.Mtime = time.Now()
	user.Ctime = time.Now()
	user.Username = req.Username
	user.Passcode = utils.RandSeq(6)
	user.Passwd = utils.Password(req.Password, user.Passcode)
	user.Hardware = req.Hardware
	_, err = db.Engine.Table(user).Insert(user)
	if err != nil {
		logger.Warn("注册数据插入失败,err: %v", err)
		return common.New(errcode.DBError, "数据库异常")
	}

	return nil
}
