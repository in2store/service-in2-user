package global

import (
	"github.com/johnnyeven/libtools/courier/transport_http"
	"github.com/johnnyeven/libtools/log"
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/libtools/sqlx/mysql"

	"github.com/in2store/service-in2-user/database"
)

func init() {
	servicex.SetServiceName("service-in2-user")
	servicex.ConfP(&Config)

	database.DBIn2Book.MustMigrateTo(Config.MasterDB.Get(), !servicex.AutoMigrate)

}

var Config = struct {
	Log    *log.Log
	Server transport_http.ServeHTTP

	MasterDB *mysql.MySQL
	SlaveDB  *mysql.MySQL
}{
	Log: &log.Log{
		Level: "DEBUG",
	},
	Server: transport_http.ServeHTTP{
		WithCORS: true,
		Port:     8000,
	},

	MasterDB: &mysql.MySQL{
		Name:     "in2-book",
		Port:     3306,
		User:     "root",
		Password: "123456",
		Host:     "127.0.0.1",
	},
	SlaveDB: &mysql.MySQL{
		Name:     "in2-book-readonly",
		Port:     3306,
		User:     "root",
		Password: "123456",
		Host:     "127.0.0.1",
	},
}
