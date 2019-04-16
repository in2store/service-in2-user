package routes

import (
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/swagger"
)

var RootRouter = courier.NewRouter(GroupRoot{})

func init() {
	RootRouter.Register(swagger.SwaggerRouter)
}

type GroupRoot struct {
	courier.EmptyOperator
}

func (root GroupRoot) Path() string {
	return "/in2-user"
}
