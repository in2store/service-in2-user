package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
)

//go:generate libtools gen model UserEntry --database DBIn2Book --table-name t_user_entry --with-comments
// @def primary ID
// @def unique_index U_user_entry_id UserEntryID
// @def unique_index U_entry EntryID ChannelID
type UserEntry struct {
	presets.PrimaryID
	// 业务ID
	UserEntryID uint64 `json:"userEntryID,string" db:"F_user_entry_id" sql:"bigint(64) unsigned NOT NULL"`
	// UserID
	UserID uint64 `json:"userID,string" db:"F_user_id" sql:"bigint(64) unsigned NOT NULL"`
	// 入口系统唯一标识
	EntryID string `json:"entryID" db:"F_entry_id" sql:"varchar(64) NOT NULL"`
	// 入口系统通道ID
	ChannelID uint64 `json:"channelID,string" db:"F_channel_id" sql:"bigint(64) unsigned NOT NULL"`

	presets.OperateTime
	presets.SoftDelete
}
