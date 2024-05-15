package core

import (
	"encoding/json"
	"os"
)

type Table struct {
	Name        string
	Description string
	Db          string
}

func (db DB) CreateTable(name string, description string) (Table, error) {
	handlers, err := db.GetTableAll()

	if err != nil {
		return Table{}, err
	}

	// 若table存在，返回存在的table
	for _, handler := range handlers {
		if handler.Name == name {
			return handler, nil
		}
	}

	// 若table不存在，创建table
	err = os.MkdirAll(DATAROOT+"/"+db.Name+"/"+name, os.ModePerm)
	if err != nil {
		return Table{}, err
	}

	// 并更新table元数据文件
	newHandler := Table{
		Name:        name,
		Description: description,
		Db:          db.Name,
	}
	tableMetadata := append(handlers, newHandler)

	err = db.updateTableMeta(tableMetadata)
	if err != nil {
		return Table{}, err
	}

	return newHandler, nil
}

func (db DB) RemoveTable(name string) error {
	handlers, err := db.GetTableAll()
	if err != nil {
		return err
	}

	err = os.RemoveAll(DATAROOT + "/" + db.Name + "/" + name)
	if err != nil {
		return err
	}

	for index, handler := range handlers {
		if handler.Name == name {
			handlers = append(handlers[:index], handlers[index+1:]...)
			db.updateTableMeta(handlers)
			break
		}
	}

	return nil
}

func (db DB) GetTableAll() ([]Table, error) {
	// 判断table元数据文件是否存在
	tableMetaPath := DATAROOT + "/" + db.Name + "/table.json"
	if _, err := os.Stat(tableMetaPath); os.IsNotExist(err) {
		// 若table元数据文件不存在，创建并返回数据
		file, err := os.Create(tableMetaPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		tableList := make([]Table, 0)
		jsonData, err := json.Marshal(tableList)

		if err != nil {
			return nil, err
		}

		file.Write(jsonData)
		return tableList, nil
	}

	// 若table元数据文件存在，返回数据
	var tableList []Table
	data, err := os.ReadFile(tableMetaPath)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &tableList)
	return tableList, nil
}

func (db DB) GetTable(name string) (Table, error) {
	handlers, err := db.GetTableAll()
	if err != nil {
		return Table{}, err
	}

	for _, handler := range handlers {
		if handler.Name == name {
			return handler, nil
		}
	}

	return Table{}, nil
}

func (db DB) updateTableMeta(data []Table) error {
	tableMetaPath := DATAROOT + "/" + db.Name + "/table.json"
	file, err := os.OpenFile(tableMetaPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(jsonData))

	if err != nil {
		return err
	}

	return nil
}
