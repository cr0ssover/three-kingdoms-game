package controller

import (
	"log"
	"net/http"

	"github.com/cr0ssover/three-kingdoms-game/server/constant/errcode"
	"github.com/cr0ssover/three-kingdoms-game/server/server/common"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web/logic"
	"github.com/cr0ssover/three-kingdoms-game/server/server/web/model"
	"github.com/gin-gonic/gin"
)

var DefaultAccountConterller = &AccountController{}

type AccountController struct {
}

func (*AccountController) Register(c *gin.Context) {
	req := &model.RegisterReq{}
	err := c.ShouldBind(req)
	if err != nil {
		log.Printf("参数有误，err:%v", err)
		c.JSON(http.StatusOK, common.Error(errcode.InvalidParam, "参数不合法"))
		return
	}
	err = logic.DefaultAccountLogic.Register(req)
	if err != nil {
		log.Printf("注册业务出错，err:%v", err)
		c.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(errcode.OK, nil))
}
