package entries

import (
	"context"
	"github.com/in2store/service-in2-user/global"
	"github.com/in2store/service-in2-user/modules/entry"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetEntries{}))
}

// 获取入口列表
type GetEntries struct {
	httpx.MethodGet
	// 用户ID
	UserID uint64 `name:"userID,string" in:"path"`
}

func (req GetEntries) Path() string {
	return ""
}

func (req GetEntries) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	result, err = entry.GetEntriesByUserID(req.UserID, db)
	if err != nil {
		logrus.Errorf("entry.GetEntriesByUserID err: %v, request: %+v", err, req)
	}
	return
}
