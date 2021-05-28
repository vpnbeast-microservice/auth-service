package database

import (
	"database/sql"
	"testing"
)

func TestGetDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root@/blog")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
}
