package errorcodes

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	MySQLErrorCodeDuplicateEntry   = 1062
	MySQLErrorCodeRowIsReferenced2 = 1451
)

func IsDuplicateEntry(err *mysql.MySQLError, value string) bool {
	if err == nil {
		return false
	}
	message := fmt.Sprintf("Duplicate entry '%s' for key", value)
	return err.Number == MySQLErrorCodeDuplicateEntry && strings.Contains(err.Message, message)
}
