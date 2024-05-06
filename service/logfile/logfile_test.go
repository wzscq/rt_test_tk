package logfile

import (
	"testing"
	"fmt"
	"rt_test_service/crv"
)

func _TestGetLogFileList(t *testing.T) {
	files, err := GetLogFileList("../mqtt/")
	if err != nil {
		t.Error("GetLogFileList error")
		return
	}

	fmt.Println(files)

	for _, file := range files {
		fmt.Println(file.Name, file.Size, file.CreationTime)
	}
}

func _TestGetLogFileFromDB(t *testing.T) {
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

func _TestUpdateLogFileToDB(t *testing.T) {
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

func _TestDeleteAllLogFiles(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	DeleteAllLogFiles(crvClient,"")
}

func _TestDeleteLogFileByName(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	DeleteLogFileByName("test",crvClient,"")
}

func _TestGetTestLogByTime(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	lfm := &LogFileMonitor{
		LogFilePath: "",
		CRVClient: crvClient,
	}

	cnt:=lfm.GetTestLogByTime("2023-09-09 13:00:33")
	fmt.Println(cnt)
}