package models

import (
	"fmt"

	"github.com/sunp13/dbtool"
)

type descModel struct{}

func (m *descModel) Desc(tableName string) (res []map[string]interface{}, err error) {
	sql := fmt.Sprintf(`
	desc %s
	`, tableName)

	res, err = dbtool.D.QuerySQL(sql, nil)
	return
}
