package users

import (
	"github.com/in2store/service-in2-user/routes/users/entries"
	"github.com/johnnyeven/libtools/courier"
)

var Router = courier.NewRouter(UsersGroup{})

func init() {
	Router.Register(entries.Router)
}

type UsersGroup struct {
	courier.EmptyOperator
}

func (UsersGroup) Path() string {
	return "/users"
}
