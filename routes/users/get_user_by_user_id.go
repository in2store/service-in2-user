package users

import (
	"context"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
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
	return nil, nil
}
