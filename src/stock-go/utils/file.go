package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type CommonCSVFile struct {
	Column []string
	Data   [][]string
}

func ReadCommonCSVFile(file string) (*CommonCSVFile, error) {
	if file == "" {
		return nil, fmt.Errorf("[ReadCommonCSVFile] file is nil")
	}
	rowData, err := ioutil.ReadFile(file)
	if err != nil || len(rowData) <= 0 {
		return nil, fmt.Errorf("[ReadCommonCSVFile] ioutil.ReadFile fail\n\t%s", err)
	}
	lines := bytes.Split(rowData, []byte("\n"))
	if len(lines) <= 0 {
		return nil, fmt.Errorf("[ReadCommonCSVFile] no file lines %s", file)
	}
	if bytes.Contains(lines[0], []byte("error")) {
		return nil, fmt.Errorf("[ReadCommonCSVFile] file data error %s", file)
	}
	columns := bytes.Split(lines[0], []byte(","))
	if len(columns) <= 0 {
		return nil, fmt.Errorf("[ReadCommonCSVFile] file data error %s", file)
	}
	ret := &CommonCSVFile{
		Column: []string{},
		Data:   [][]string{},
	}
	for i := 0; i < len(columns); i++ {
		ret.Column = append(ret.Column, string(columns[i]))
	}
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) <= 0 {
			// 跳过空行
			continue
		}
		rowColumn := bytes.Split(lines[i], []byte(","))
		if len(rowColumn) != len(ret.Column) {
			return nil, fmt.Errorf("[ReadCommonCSVFile] file data error %s %s", lines[i], file)
		}
		row := []string{}
		for j := 0; j < len(rowColumn); j++ {
			row = append(row, string(rowColumn[j]))
		}
		ret.Data = append(ret.Data, row)
	}
	return ret, nil
}
