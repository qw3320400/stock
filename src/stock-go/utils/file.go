package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	ErrFileNotExist = fmt.Errorf("No such file or directory")

	ErrFileQueryTimeout = fmt.Errorf("File query timeout")
)

type CommonCSVFile struct {
	Column []string
	Data   [][]string
}

func ReadCommonCSVFile(file string) (*CommonCSVFile, error) {
	if file == "" {
		return nil, Errorf(nil, "file is nil")
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, ErrFileNotExist
	}
	rowData, err := ioutil.ReadFile(file)
	if err != nil || len(rowData) <= 0 {
		return nil, Errorf(err, "ioutil.ReadFile fail")
	}
	lines := bytes.Split(rowData, []byte("\n"))
	if len(lines) <= 0 {
		return nil, Errorf(nil, "no file lines %s", file)
	}
	if bytes.Contains(lines[0], []byte("error")) {
		return nil, ErrFileQueryTimeout
	}
	columns := bytes.Split(lines[0], []byte(","))
	if len(columns) <= 0 {
		return nil, Errorf(nil, "file data error %s", file)
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
			return nil, Errorf(nil, "file data error %s %s", lines[i], file)
		}
		row := []string{}
		for j := 0; j < len(rowColumn); j++ {
			row = append(row, string(rowColumn[j]))
		}
		ret.Data = append(ret.Data, row)
	}
	return ret, nil
}
