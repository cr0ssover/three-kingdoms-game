package data

import (
	"time"

	"github.com/cr0ssover/three-kingdoms-game/server/server/game/model"
)

// 玩家属性表
type RoleAttribute struct {
	Id              int            `xorm:"int(0) pk autoincr 'id' comment('id')"`
	RId             int            `xorm:"int(0) 'rid' comment('角色id')"`
	UnionId         int            `xorm:"-"`                                                  //联盟id
	ParentId        int            `xorm:"int(0) 'parent_id' comment('上级联盟id')"`               //上级id（被沦陷）
	CollectTimes    int8           `xorm:"int(0) 'collect_times' comment('征收次数')"`             //征收次数
	LastCollectTime time.Time      `xorm:"timestamp(0) 'last_collect_time' comment('最后征收时间')"` //最后征收的时间
	PosTags         string         `xorm:"varchar(512) 'pos_tags' comment('收藏的位置')"`           //位置标记
	PosTagArray     []model.PosTag `xorm:"-"`
}

func (r *RoleAttribute) TableName() string {
	return "role_attribute"
}
