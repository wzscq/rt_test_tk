package logfile

import (
	"testing"
	"fmt"
	"rt_test_service/crv"
)

func TestGetLogFileList(t *testing.T) {
	files, err := GetLogFileList(".")
	if err != nil {
		t.Error("GetLogFileList error")
		return
	}

	for _, file := range files {
		fmt.Println(file.Name, file.Size, file.CreationTime)
	}
}

func TestGetLogFileFromDB(t *testing.T) {
	file := LogFileItem{
		Name: "test",
		Size: 100,
		CreationTime: "2019-01-01 12:00:00",
	}

	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	fileMap,err := GetLogFileFromDB(file, crvClient, "rt_test_tk_service")
	if err != nil {
		t.Error("GetLogFileFromDB error")
		return
	}

	if fileMap == nil {
		t.Error("GetLogFileFromDB error")
		return
	}

	fmt.Println(fileMap)
}

func TestUpdateLogFileToDB(t *testing.T) {
	files, err := GetLogFileList(".")
	if err != nil {
		t.Error("GetLogFileList error")
		return
	}

	if len(files) == 0 {
		t.Error("GetLogFileList error")
		return
	}

	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	err=UpdateLogFilesToDB(files, crvClient, "rt_test_tk_service")
	if err != nil {
		t.Error("UpdateLogFilesToDB error")
		return
	}
}