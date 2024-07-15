package errorcodes

import (
	"testing"

	"github.com/go-sql-driver/mysql"
)

func Test_IsDuplicateEntry_NilError_False(t *testing.T) {
	if IsDuplicateEntry(nil, "value") {
		t.Error("Expected false")
	}
}

func Test_IsDuplicateEntry_NotDuplicateEntry_False(t *testing.T) {
	err := &mysql.MySQLError{Number: 1, Message: "message"}
	if IsDuplicateEntry(err, "value") {
		t.Error("Expected false")
	}
}

func Test_IsDuplicateEntry_DuplicateEntry_False(t *testing.T) {
	err := &mysql.MySQLError{Number: MySQLErrorCodeDuplicateEntry, Message: "message"}
	if IsDuplicateEntry(err, "value") {
		t.Error("Expected false")
	}
}

func Test_IsDuplicateEntry_DuplicateEntry_True(t *testing.T) {
	err := &mysql.MySQLError{Number: MySQLErrorCodeDuplicateEntry, Message: "Duplicate entry 'value' for key"}
	if !IsDuplicateEntry(err, "value") {
		t.Error("Expected true")
	}
}
