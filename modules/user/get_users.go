package user

import (
	"github.com/in2store/service-in2-user/database"
	"github.com/johnnyeven/libtools/httplib"
	"github.com/johnnyeven/libtools/sqlx"
	"github.com/johnnyeven/libtools/sqlx/builder"
)

type GetUsersParams struct {
	// 入口ID
	EntryID string
	// 通道ID
	ChannelID uint64
	// 用户ID列表
	UserIDs httplib.Uint64List
	// 分页
	Size int32
	// 偏移量
	Offset int32
}

func (p GetUsersParams) Conditions() *builder.Condition {
	t := (&database.User{}).T()
	var condition *builder.Condition
	if len(p.UserIDs) > 0 {
		condition = builder.And(condition, t.F("UserID").In(p.UserIDs))
	}

	return condition
}

func GetUsers(params GetUsersParams, db *sqlx.DB) (result database.UserList, count int32, err error) {
	if params.EntryID != "" || params.ChannelID != 0 {
		entry := &database.UserEntry{}
		entries, _, err := entry.FetchList(db, -1, 0)
		if err != nil {
			return nil, 0, err
		}

		for _, e := range entries {
			params.UserIDs = append(params.UserIDs, e.UserID)
		}
	}

	user := &database.User{}
	result, count, err = user.FetchList(db, params.Size, params.Offset, params.Conditions())
	return
}

func GetUserByUserID(userID uint64, db *sqlx.DB) (*database.User, error) {
	user := &database.User{
		UserID: userID,
	}
	err := user.FetchByUserID(db)
	if err != nil {
		return nil, err
	}
	return user, nil
}
