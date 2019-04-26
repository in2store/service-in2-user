package entry

import (
	"github.com/in2store/service-in2-user/database"
	"github.com/johnnyeven/libtools/sqlx"
)

func GetEntriesByUserID(userID uint64, db *sqlx.DB) (database.UserEntryList, error) {
	entry := &database.UserEntry{}
	entries, err := entry.BatchFetchByUserIDList(db, []uint64{userID})
	if err != nil {
		return nil, err
	}
	return entries, nil
}
