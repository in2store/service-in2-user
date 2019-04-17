package users

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(UsersGroup{})

type UsersGroup struct {
	courier.EmptyOperator
}

func (UsersGroup) Path() string {
	return "/users"
}
