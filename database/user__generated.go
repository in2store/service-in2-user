package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var UserTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	UserTable = DBIn2Book.Register(&User{})
}

func (user *User) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (user *User) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return UserTable
}

func (user *User) TableName() string {
	return "t_user"
}

type UserFields struct {
	ID         *github_com_johnnyeven_libtools_sqlx_builder.Column
	UserID     *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled    *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var UserField = struct {
	ID         string
	UserID     string
	CreateTime string
	UpdateTime string
	Enabled    string
}{
	ID:         "ID",
	UserID:     "UserID",
	CreateTime: "CreateTime",
	UpdateTime: "UpdateTime",
	Enabled:    "Enabled",
}

func (user *User) Fields() *UserFields {
	table := user.T()

	return &UserFields{
		ID:         table.F(UserField.ID),
		UserID:     table.F(UserField.UserID),
		CreateTime: table.F(UserField.CreateTime),
		UpdateTime: table.F(UserField.UpdateTime),
		Enabled:    table.F(UserField.Enabled),
	}
}

func (user *User) IndexFieldNames() []string {
	return []string{"ID", "UserID"}
}

func (user *User) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := user.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(user)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range user.IndexFieldNames() {
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

func (user *User) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (user *User) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{"U_user_id": github_com_johnnyeven_libtools_sqlx.FieldNames{"UserID", "Enabled"}}
}
func (user *User) Comments() map[string]string {
	return map[string]string{
		"CreateTime": "",
		"Enabled":    "",
		"ID":         "",
		"UpdateTime": "",
		"UserID":     "业务ID",
	}
}

func (user *User) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if user.CreateTime.IsZero() {
		user.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	user.UpdateTime = user.CreateTime

	stmt := user.D().
		Insert(user).
		Comment("User.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		user.ID = uint64(lastInsertID)
	}

	return err
}

func (user *User) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := user.T()

	stmt := table.Delete().
		Comment("User.DeleteByStruct").
		Where(user.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (user *User) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if user.CreateTime.IsZero() {
		user.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	user.UpdateTime = user.CreateTime

	table := user.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(user, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range user.UniqueIndexes() {
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
		Comment("User.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (user *User) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Delete().
		Comment("User.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("User.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return user.FetchByID(db)
	}
	return nil
}

func (user *User) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(user, zeroFields...)
	return user.UpdateByIDWithMap(db, fieldValues)
}

func (user *User) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("User.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(user.ID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return user.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (user *User) FetchByUserID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByUserID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(user.UserID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) FetchByUserIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Select().
		Comment("User.FetchByUserIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(user.UserID),
			table.F("Enabled").Eq(user.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) DeleteByUserID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()
	stmt := table.Delete().
		Comment("User.DeleteByUserID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(user.UserID),
			table.F("Enabled").Eq(user.Enabled),
		))

	return db.Do(stmt).Scan(user).Err()
}

func (user *User) UpdateByUserIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("User.UpdateByUserIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(user.UserID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return user.FetchByUserID(db)
	}
	return nil
}

func (user *User) UpdateByUserIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(user, zeroFields...)
	return user.UpdateByUserIDWithMap(db, fieldValues)
}

func (user *User) SoftDeleteByUserID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	user.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := user.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("User.SoftDeleteByUserID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(user.UserID),
			table.F("Enabled").Eq(user.Enabled),
		))

	dbRet := db.Do(stmt).Scan(user)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return user.DeleteByUserID(db)
		}
		return err
	}
	return nil
}

type UserList []User

// deprecated
func (userList *UserList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*userList, count, err = (&User{}).FetchList(db, size, offset, conditions...)
	return
}

func (user *User) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (userList UserList, count int32, err error) {
	userList = UserList{}

	table := user.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

func (user *User) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (userList UserList, err error) {
	userList = UserList{}

	table := user.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.List").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

func (user *User) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (userList UserList, err error) {
	userList = UserList{}

	table := user.T()

	condition := user.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*userList, err = (&User{}).BatchFetchByIDList(db, idList)
	return
}

func (user *User) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (userList UserList, err error) {
	if len(idList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}

// deprecated
func (userList *UserList) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (err error) {
	*userList, err = (&User{}).BatchFetchByUserIDList(db, userIDList)
	return
}

func (user *User) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (userList UserList, err error) {
	if len(userIDList) == 0 {
		return UserList{}, nil
	}

	table := user.T()

	condition := table.F("UserID").In(userIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("User.BatchFetchByUserIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&userList).Err()

	return
}
