package users

import (
	"context"
	"github.com/in2store/service-in2-user/global"
	"github.com/in2store/service-in2-user/modules/user"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetUserByUserID{}))
}

// 根据用户ID获取用户详情
type GetUserByUserID struct {
	httpx.MethodGet
	// 用户ID
	UserID uint64 `name:"userID,string" in:"path"`
}

func (req GetUserByUserID) Path() string {
	return "/:userID"
}

func (req GetUserByUserID) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	u, err := user.GetUserByUserID(req.UserID, db)
	if err != nil {
		logrus.Errorf("[GetUserByUserID] user.GetUserByUserID err: %v, request: %+v", err, req)
		return nil, err
	}
	return u, nil
}
