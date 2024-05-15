package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

type FieldType = uint32

const (
	INT FieldType = iota
	UINT
	FLOAT
	STRING
	BOOL
)

type FieldMeta struct {
	Name         string
	Description  string
	Type         FieldType
	Required     bool
	Readonly     bool
	DefaultValue interface{}
}

func (table Table) DefineFieldMeta(metadata []FieldMeta) error {
	metaPath := DATAROOT + "/" + table.Db + "/" + table.Name + "/meta.json"
	var file *os.File

	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		// 若文件不存在，创建并返回数据
		file, err = os.Create(metaPath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(metaPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
	}

	defer file.Close()

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(jsonData))

	if err != nil {
		return err
	}

	return nil
}

func (table Table) GetFieldMeta() ([]FieldMeta, error) {
	metaPath := DATAROOT + "/" + table.Db + "/" + table.Name + "/meta.json"
	var metaData []FieldMeta
	data, err := os.ReadFile(metaPath)

	if err != nil {
		return []FieldMeta{}, err
	}

	json.Unmarshal(data, &metaData)
	return metaData, nil
}

func (table Table) Create(dataList []interface{}) error {
	dataPath := DATAROOT + "/" + table.Db + "/" + table.Name + "/data.json"
	var file *os.File

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		file, err = os.Create(dataPath)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
	}

	defer file.Close()

	metadata, err := table.GetFieldMeta()
	if err != nil {
		return err
	}

	for _, data := range dataList {
		err := validateData(data, metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateData(data interface{}, metadata []FieldMeta) error {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("%s : %v", f.Name, f.Type)
		fmt.Println("val :", val)

		_, err := getMetaData(metadata, f.Name)
		if err != nil {
			return err
		}

	}

	return nil
}

func getMetaData(metadata []FieldMeta, name string) (FieldMeta, error) {
	for _, data := range metadata {
		if data.Name == name {
			return data, nil
		}
	}

	return FieldMeta{}, errors.New("field not found")
}
