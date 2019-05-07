package database

import (
	"github.com/johnnyeven/libtools/sqlx/presets"
)

//go:generate libtools gen model User --database DBIn2Book --table-name t_user --with-comments
// @def primary ID
// @def unique_index U_user_id UserID
type User struct {
	presets.PrimaryID
	// 业务ID
	UserID uint64 `json:"userID,string" db:"F_user_id" sql:"bigint(64) unsigned NOT NULL"`
	// 用户名
	Name string `json:"name" db:"F_name" sql:"varchar(32) NOT NULL"`

	presets.OperateTime
	presets.SoftDelete
}
