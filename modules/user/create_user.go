package user

import (
	"github.com/in2store/service-in2-user/constants/errors"
	"github.com/in2store/service-in2-user/database"
	"github.com/johnnyeven/eden-library/clients/client_id"
	"github.com/johnnyeven/eden-library/modules"
	"github.com/johnnyeven/libtools/sqlx"
)

type CreateUserParams struct {
	Entries []CreateUserParamsEntry `json:"entries"`
}

func (req CreateUserParams) Validate() error {
	if len(req.Entries) == 0 {
		return errors.BadRequest.StatusError().WithMsg("入口信息不能为空").WithErrTalk()
	}
	return nil
}

type CreateUserParamsEntry struct {
	// 入口系统唯一标识
	EntryID string `json:"entryID"`
	// 入口系统通道ID
	ChannelID uint64 `json:"channelID,string"`
}

func CreateUser(req CreateUserParams, db *sqlx.DB, client *client_id.ClientID) (user *database.User, err error) {
	tx := sqlx.NewTasks(db)
	var userID uint64
	createUser := func(db *sqlx.DB) error {
		userID, err = modules.NewUniqueID(client)
		if err != nil {
			return err
		}
		user = &database.User{
			UserID: userID,
		}
		return user.Create(db)
	}
	tx = tx.With(createUser)

	for _, entry := range req.Entries {
		createEntry := func(db *sqlx.DB) error {
			userEntryID, err := modules.NewUniqueID(client)
			if err != nil {
				return err
			}
			e := &database.UserEntry{
				UserEntryID: userEntryID,
				UserID:      userID,
				EntryID:     entry.EntryID,
				ChannelID:   entry.ChannelID,
			}
			return e.Create(db)
		}
		tx = tx.With(createEntry)
	}

	err = tx.Do()
	if err != nil {
		return nil, err
	}

	return user, nil
}
