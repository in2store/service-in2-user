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
	Router.Register(courier.NewRouter(CreateUser{}))
}

// 创建用户
type CreateUser struct {
	httpx.MethodPost
	Body user.CreateUserParams `name:"body" in:"body"`
}

func (req CreateUser) Path() string {
	return ""
}

func (req CreateUser) Output(ctx context.Context) (result interface{}, err error) {
	err = req.Body.Validate()
	if err != nil {
		logrus.Warningf("[CreateUser] 参数校验异常 err: %v, request: %+v", err, req.Body)
		return nil, err
	}

	db := global.Config.MasterDB.Get()
	result, err = user.CreateUser(req.Body, db, global.Config.ClientID)
	if err != nil {
		logrus.Errorf("[CreateUser] user.CreateUser err: %v, request: %+v", err, req.Body)
		return nil, err
	}

	return
}
