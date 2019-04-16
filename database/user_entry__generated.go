package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var UserEntryTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	UserEntryTable = DBIn2Book.Register(&UserEntry{})
}

func (userEntry *UserEntry) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (userEntry *UserEntry) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return UserEntryTable
}

func (userEntry *UserEntry) TableName() string {
	return "t_user_entry"
}

type UserEntryFields struct {
	ID          *github_com_johnnyeven_libtools_sqlx_builder.Column
	UserEntryID *github_com_johnnyeven_libtools_sqlx_builder.Column
	UserID      *github_com_johnnyeven_libtools_sqlx_builder.Column
	EntryID     *github_com_johnnyeven_libtools_sqlx_builder.Column
	ChannelID   *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime  *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime  *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled     *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var UserEntryField = struct {
	ID          string
	UserEntryID string
	UserID      string
	EntryID     string
	ChannelID   string
	CreateTime  string
	UpdateTime  string
	Enabled     string
}{
	ID:          "ID",
	UserEntryID: "UserEntryID",
	UserID:      "UserID",
	EntryID:     "EntryID",
	ChannelID:   "ChannelID",
	CreateTime:  "CreateTime",
	UpdateTime:  "UpdateTime",
	Enabled:     "Enabled",
}

func (userEntry *UserEntry) Fields() *UserEntryFields {
	table := userEntry.T()

	return &UserEntryFields{
		ID:          table.F(UserEntryField.ID),
		UserEntryID: table.F(UserEntryField.UserEntryID),
		UserID:      table.F(UserEntryField.UserID),
		EntryID:     table.F(UserEntryField.EntryID),
		ChannelID:   table.F(UserEntryField.ChannelID),
		CreateTime:  table.F(UserEntryField.CreateTime),
		UpdateTime:  table.F(UserEntryField.UpdateTime),
		Enabled:     table.F(UserEntryField.Enabled),
	}
}

func (userEntry *UserEntry) IndexFieldNames() []string {
	return []string{"ChannelID", "EntryID", "ID", "UserEntryID"}
}

func (userEntry *UserEntry) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := userEntry.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(userEntry)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range userEntry.IndexFieldNames() {
		if v, exists := fieldValues[fieldName]; exists {
			conditions = append(conditions, table.F(fieldName).Eq(v))
			delete(fieldValues, fieldName)
		}
	}

	if len(conditions) == 0 {
		panic(fmt.Errorf("at least one of field for indexes has value"))
	}

	for fieldName, v := range fieldValues {
		conditions = append(conditions, table.F(fieldName).Eq(v))
	}

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	return condition
}

func (userEntry *UserEntry) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (userEntry *UserEntry) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{
		"U_entry":         github_com_johnnyeven_libtools_sqlx.FieldNames{"EntryID", "ChannelID", "Enabled"},
		"U_user_entry_id": github_com_johnnyeven_libtools_sqlx.FieldNames{"UserEntryID", "Enabled"},
	}
}
func (userEntry *UserEntry) Comments() map[string]string {
	return map[string]string{
		"ChannelID":   "入口系统通道ID",
		"CreateTime":  "",
		"Enabled":     "",
		"EntryID":     "入口系统唯一标识",
		"ID":          "",
		"UpdateTime":  "",
		"UserEntryID": "业务ID",
		"UserID":      "UserID",
	}
}

func (userEntry *UserEntry) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if userEntry.CreateTime.IsZero() {
		userEntry.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	userEntry.UpdateTime = userEntry.CreateTime

	stmt := userEntry.D().
		Insert(userEntry).
		Comment("UserEntry.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		userEntry.ID = uint64(lastInsertID)
	}

	return err
}

func (userEntry *UserEntry) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := userEntry.T()

	stmt := table.Delete().
		Comment("UserEntry.DeleteByStruct").
		Where(userEntry.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (userEntry *UserEntry) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if userEntry.CreateTime.IsZero() {
		userEntry.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	userEntry.UpdateTime = userEntry.CreateTime

	table := userEntry.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(userEntry, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range userEntry.UniqueIndexes() {
		for _, field := range fieldNames {
			delete(m, field)
		}
	}

	if len(m) == 0 {
		panic(fmt.Errorf("no fields for updates"))
	}

	for field := range fieldValues {
		if !m[field] {
			delete(fieldValues, field)
		}
	}

	stmt := table.
		Insert().Columns(cols).Values(vals...).
		OnDuplicateKeyUpdate(table.AssignsByFieldValues(fieldValues)...).
		Comment("UserEntry.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (userEntry *UserEntry) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(userEntry.ID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(userEntry.ID),
			table.F("Enabled").Eq(userEntry.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Delete().
		Comment("UserEntry.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(userEntry.ID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("UserEntry.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(userEntry.ID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return userEntry.FetchByID(db)
	}
	return nil
}

func (userEntry *UserEntry) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(userEntry, zeroFields...)
	return userEntry.UpdateByIDWithMap(db, fieldValues)
}

func (userEntry *UserEntry) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("UserEntry.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(userEntry.ID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return userEntry.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (userEntry *UserEntry) FetchByUserEntryID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByUserEntryID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserEntryID").Eq(userEntry.UserEntryID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) FetchByUserEntryIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByUserEntryIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserEntryID").Eq(userEntry.UserEntryID),
			table.F("Enabled").Eq(userEntry.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) DeleteByUserEntryID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Delete().
		Comment("UserEntry.DeleteByUserEntryID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserEntryID").Eq(userEntry.UserEntryID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) UpdateByUserEntryIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("UserEntry.UpdateByUserEntryIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserEntryID").Eq(userEntry.UserEntryID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return userEntry.FetchByUserEntryID(db)
	}
	return nil
}

func (userEntry *UserEntry) UpdateByUserEntryIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(userEntry, zeroFields...)
	return userEntry.UpdateByUserEntryIDWithMap(db, fieldValues)
}

func (userEntry *UserEntry) SoftDeleteByUserEntryID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("UserEntry.SoftDeleteByUserEntryID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserEntryID").Eq(userEntry.UserEntryID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return userEntry.DeleteByUserEntryID(db)
		}
		return err
	}
	return nil
}

func (userEntry *UserEntry) FetchByEntryIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByEntryIDAndChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("EntryID").Eq(userEntry.EntryID),
			table.F("ChannelID").Eq(userEntry.ChannelID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) FetchByEntryIDAndChannelIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Select().
		Comment("UserEntry.FetchByEntryIDAndChannelIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("EntryID").Eq(userEntry.EntryID),
			table.F("ChannelID").Eq(userEntry.ChannelID),
			table.F("Enabled").Eq(userEntry.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) DeleteByEntryIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()
	stmt := table.Delete().
		Comment("UserEntry.DeleteByEntryIDAndChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("EntryID").Eq(userEntry.EntryID),
			table.F("ChannelID").Eq(userEntry.ChannelID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	return db.Do(stmt).Scan(userEntry).Err()
}

func (userEntry *UserEntry) UpdateByEntryIDAndChannelIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("UserEntry.UpdateByEntryIDAndChannelIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("EntryID").Eq(userEntry.EntryID),
			table.F("ChannelID").Eq(userEntry.ChannelID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return userEntry.FetchByEntryIDAndChannelID(db)
	}
	return nil
}

func (userEntry *UserEntry) UpdateByEntryIDAndChannelIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(userEntry, zeroFields...)
	return userEntry.UpdateByEntryIDAndChannelIDWithMap(db, fieldValues)
}

func (userEntry *UserEntry) SoftDeleteByEntryIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	userEntry.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := userEntry.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("UserEntry.SoftDeleteByEntryIDAndChannelID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("EntryID").Eq(userEntry.EntryID),
			table.F("ChannelID").Eq(userEntry.ChannelID),
			table.F("Enabled").Eq(userEntry.Enabled),
		))

	dbRet := db.Do(stmt).Scan(userEntry)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return userEntry.DeleteByEntryIDAndChannelID(db)
		}
		return err
	}
	return nil
}

type UserEntryList []UserEntry

// deprecated
func (userEntryList *UserEntryList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*userEntryList, count, err = (&UserEntry{}).FetchList(db, size, offset, conditions...)
	return
}

func (userEntry *UserEntry) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (userEntryList UserEntryList, count int32, err error) {
	userEntryList = UserEntryList{}

	table := userEntry.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

func (userEntry *UserEntry) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (userEntryList UserEntryList, err error) {
	userEntryList = UserEntryList{}

	table := userEntry.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.List").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

func (userEntry *UserEntry) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (userEntryList UserEntryList, err error) {
	userEntryList = UserEntryList{}

	table := userEntry.T()

	condition := userEntry.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

// deprecated
func (userEntryList *UserEntryList) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (err error) {
	*userEntryList, err = (&UserEntry{}).BatchFetchByChannelIDList(db, channelIDList)
	return
}

func (userEntry *UserEntry) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (userEntryList UserEntryList, err error) {
	if len(channelIDList) == 0 {
		return UserEntryList{}, nil
	}

	table := userEntry.T()

	condition := table.F("ChannelID").In(channelIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.BatchFetchByChannelIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

// deprecated
func (userEntryList *UserEntryList) BatchFetchByEntryIDList(db *github_com_johnnyeven_libtools_sqlx.DB, entryIDList []string) (err error) {
	*userEntryList, err = (&UserEntry{}).BatchFetchByEntryIDList(db, entryIDList)
	return
}

func (userEntry *UserEntry) BatchFetchByEntryIDList(db *github_com_johnnyeven_libtools_sqlx.DB, entryIDList []string) (userEntryList UserEntryList, err error) {
	if len(entryIDList) == 0 {
		return UserEntryList{}, nil
	}

	table := userEntry.T()

	condition := table.F("EntryID").In(entryIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.BatchFetchByEntryIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

// deprecated
func (userEntryList *UserEntryList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*userEntryList, err = (&UserEntry{}).BatchFetchByIDList(db, idList)
	return
}

func (userEntry *UserEntry) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (userEntryList UserEntryList, err error) {
	if len(idList) == 0 {
		return UserEntryList{}, nil
	}

	table := userEntry.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}

// deprecated
func (userEntryList *UserEntryList) BatchFetchByUserEntryIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userEntryIDList []uint64) (err error) {
	*userEntryList, err = (&UserEntry{}).BatchFetchByUserEntryIDList(db, userEntryIDList)
	return
}

func (userEntry *UserEntry) BatchFetchByUserEntryIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userEntryIDList []uint64) (userEntryList UserEntryList, err error) {
	if len(userEntryIDList) == 0 {
		return UserEntryList{}, nil
	}

	table := userEntry.T()

	condition := table.F("UserEntryID").In(userEntryIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("UserEntry.BatchFetchByUserEntryIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userEntryList).Err()

	return
}
