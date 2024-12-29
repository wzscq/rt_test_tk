package logfile

import (
	"testing"
	"fmt"
	"rt_test_service/crv"
	"time"
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

	err=UpdateLogFilesToDB(files, crvClient, "rt_test_tk_service", time.Second*1800)
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

	duration, _ := time.ParseDuration("1800s")

	lfm := &LogFileMonitor{
		LogFilePath: "",
		CRVClient: crvClient,
		ExpandTimeRange: duration,
	}

	cnt:=GetTestLogByTime("2023-09-09 13:00:33",lfm.CRVClient, "rt_test_tk_service", lfm.ExpandTimeRange)
	fmt.Println(cnt)
}

func TestGetCreateTime(t *testing.T) {
	tm:=GetCreateTime("20240519_073626_0014.qmdl2")
	if tm != "2024-05-19 15:36:26" {
		t.Error("GetCreateTime error")
	}
	
	fmt.Println(tm)
}