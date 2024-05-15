package core

import (
	"encoding/json"
	"os"
)

const DATAROOT = ".data"
const DBMETAPATH = DATAROOT + "/" + "db.json"

type DB struct {
	Name        string
	Description string
}

func CreateDB(name string, description string) (DB, error) {
	handlers, err := GetDBAll()

	if err != nil {
		return DB{}, err
	}

	// 若db存在，返回存在的db
	for _, handler := range handlers {
		if handler.Name == name {
			return handler, nil
		}
	}

	// 若db不存在，创建db
	err = os.MkdirAll(DATAROOT+"/"+name, os.ModePerm)
	if err != nil {
		return DB{}, err
	}

	// 并更新db元数据文件
	newHandler := DB{
		Name:        name,
		Description: description,
	}
	dbMetadata := append(handlers, newHandler)

	err = updateDBMeta(dbMetadata)
	if err != nil {
		return DB{}, err
	}

	return newHandler, nil
}

func RemoveDB(name string) error {
	handlers, err := GetDBAll()
	if err != nil {
		return err
	}

	err = os.RemoveAll(DATAROOT + "/" + name)
	if err != nil {
		return err
	}

	for index, handler := range handlers {
		if handler.Name == name {
			handlers = append(handlers[:index], handlers[index+1:]...)
			updateDBMeta(handlers)
			break
		}
	}

	return nil
}

func GetDBAll() ([]DB, error) {
	// 判断db元数据文件是否存在
	if _, err := os.Stat(DBMETAPATH); os.IsNotExist(err) {
		// 若db元数据文件不存在，创建并返回数据
		file, err := os.Create(DBMETAPATH)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		dbList := make([]DB, 0)
		jsonData, err := json.Marshal(dbList)

		if err != nil {
			return nil, err
		}

		file.Write(jsonData)
		return dbList, nil
	}

	// 若db元数据文件存在，返回数据
	var dbList []DB
	data, err := os.ReadFile(DBMETAPATH)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &dbList)
	return dbList, nil
}

func GetDB(name string) (DB, error) {
	handlers, err := GetDBAll()
	if err != nil {
		return DB{}, err
	}

	for _, handler := range handlers {
		if handler.Name == name {
			return handler, nil
		}
	}

	return DB{}, nil
}

func updateDBMeta(data []DB) error {
	file, err := os.OpenFile(DBMETAPATH, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
