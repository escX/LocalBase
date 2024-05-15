package main

import (
	"LocalBase/core"
)

func main() {
	db, _ := core.CreateDB("testDb", "testDb")
	table, _ := db.CreateTable("test_table", "test_table")
	table.DefineFieldMeta([]core.FieldMeta{
		{"name", "name", 0, true, false, nil},
		{"age", "age", 1, true, false, nil},
	})

	type User struct {
		Name string
		Age  int
	}

	table.Create([]interface{}{
		User{
			Name: "test",
			Age:  18,
		},
	})
}
