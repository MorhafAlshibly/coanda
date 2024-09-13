package errorcode

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	MySQLErrorCodeDuplicateEntry   = 1062
	MySQLErrorCodeRowIsReferenced2 = 1451
	MySQLErrorCodeNoReferencedRow2 = 1452
)

func IsDuplicateEntry(err *mysql.MySQLError, tableName string, constraintName string) bool {
	if err == nil {
		return false
	}
	suffix := fmt.Sprintf("%s.%s", tableName, constraintName)
	return err.Number == MySQLErrorCodeDuplicateEntry && strings.Contains(err.Message, suffix)
}
