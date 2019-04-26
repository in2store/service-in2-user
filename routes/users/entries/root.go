package entries

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(EntriesGroup{})

type EntriesGroup struct {
	courier.EmptyOperator
}

func (EntriesGroup) Path() string {
	return "/:userID/entries"
}
