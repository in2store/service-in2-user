package users

import (
	"context"
	"github.com/in2store/service-in2-user/database"
	"github.com/in2store/service-in2-user/global"
	"github.com/in2store/service-in2-user/modules/user"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/httpx"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.Register(courier.NewRouter(GetUsers{}))
}

// 根据入口ID与通道ID获取用户详情
type GetUsers struct {
	httpx.MethodGet
	// 入口ID
	EntryID string `name:"entryID" in:"query"`
	// 通道ID
	ChannelID uint64 `name:"channelID,string" in:"query"`
	// 分页
	Size int32 `name:"size" in:"query" default:"10" validate:"@int32[-1,100]"`
	// 偏移量
	Offset int32 `name:"offset" in:"query" default:"0" validate:"@int32[0,]"`
}

func (req GetUsers) Path() string {
	return ""
}

type GetUsersResponse struct {
	Data  database.UserList `json:"data"`
	Total int32             `json:"total"`
}

func (req GetUsers) Output(ctx context.Context) (result interface{}, err error) {
	db := global.Config.SlaveDB.Get()
	request := user.GetUsersParams{
		EntryID:   req.EntryID,
		ChannelID: req.ChannelID,
		Size:      req.Size,
		Offset:    req.Offset,
	}
	resp, count, err := user.GetUsers(request, db)
	if err != nil {
		logrus.Errorf("[GetUsers] user.GetUsers err: %v, request: %+v", err, request)
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	return GetUsersResponse{
		Data:  resp,
		Total: count,
	}, nil
}
