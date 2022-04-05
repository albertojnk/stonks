package ports

import "database/sql"

type Persistance interface {
	Exec(string, ...interface{}) (bool, error)
	QueryRow(string, ...interface{}) (*sql.Rows, error)
	Shutdown()
	RunFileQuery(string) error
}
