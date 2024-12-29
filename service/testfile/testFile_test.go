package testfile

import (
	"encoding/json"
	"fmt"
	"log"
	"rt_test_service/crv"
	"rt_test_service/common"
	"strconv"
	"testing"
	"time"
)

var data1=map[string]interface{}{
	"exampleCode": "iperf_executor",  
	"msg_type": "report", 
	"testData": map[string]interface{}{
		"gps": map[string]interface{}{
			"gps_status": "Valid", 
			"gps_ant_status": "OK", 
			"datetime": "2024-04-05 15:59:24.473338", 
			"longitude_direction": "E", 
			"longitude": 118.8350601196289, 
			"latitude_direction": "N", 
			"latitude": 31.9070987701416,
		},
	}, 
}

var data2= map[string]interface{}{
	"exampleCode": "iperf_executor",  
	"msg_type": "finally", 
	"testData": map[string]interface{}{
		"result": "success",
		"gps": map[string]interface{}{
			"gps_status": "Valid", 
			"gps_ant_status": "OK", 
			"datetime": "2024-04-05 15:59:24.473338", 
			"longitude_direction": "E", 
			"longitude": 118.8350601196289, 
			"latitude_direction": "N", 
			"latitude": 31.9070987701416,
		},
	}, 
}

func _TestCreateFile(t *testing.T) {
	outPath := "../localcache/"
	deviceID := "device1"
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	tf := GetTestFile(outPath, deviceID, timeStamp)
	if tf == nil {
		t.Error("GetTestFile failed")
		return
	}

	for i := 0; i < 10; i++ {
		lineContent := "line" + strconv.Itoa(i)
		tf.WriteLine(lineContent)
	}

	tf.Close()
}

func _TestGetTestFileFromDB(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:  "rt_test_tk_service",
		AppID:  "rt_test_tk",
	}

	tf, errorCode := GetTestFileFromDB("32","",crvClient)
	if errorCode != common.ResultSuccess {
		t.Error("GetTestFileFromDB failed")
		return
	}

	log.Println(tf)
}

func _TestTestFilePool(t *testing.T) {

	crvClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:  "rt_test_tk_service",
		AppID:  "rt_test_tk",
	}

	testCmdSender:=&TestCmdSender{}

	//reportData转换为JSON字符串
	reportDataJson, _ := json.Marshal(data1)
	//创建TestFilePool
	tfp := InitTestFilePool("../localcache/", "3s", crvClient)
	tfp.SetCmdSender(testCmdSender)
	//写入数据
	tfp.HandleReportResult(string(reportDataJson))
	//等待5秒
	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)
		tfp.HandleReportResult(string(reportDataJson))
	}
	result,_:=json.Marshal(data2)
	tfp.HandleReportResult(string(result))
	time.Sleep(2 * time.Second)
	tfp.HandleReportResult(string(result))
	time.Sleep(2 * time.Second)
	tfp.HandleReportResult(string(result))

	tfp.SetCmdSender(nil)
	//tfp.HandleReportResult(string(result))
	time.Sleep(8 * time.Second)
}

func _TestGetFilePoints(t *testing.T) {
	outPath := "../localcache"
	deviceID := "BFEBFBFF000806EC"
	timestamp := "1694409161"

	tf := GetTestFile(outPath, deviceID, timestamp)
	if tf == nil {
		t.Error("TestFileController GetPoints file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	indicator := Indicator{
		ExtractPath: "radio.measures_lte.RSRP",
		ID:          "1",
		Name:        "test",
		Legend: IndicatorLegend{
			ModelID: "model",
			Total:   3,
			List: []IndicatorLegendItem{
				IndicatorLegendItem{
					ID:    "id",
					SN:    "1",
					Start: "",
					End:   "-80.75",
					RGB:   "#110000",
				},
				IndicatorLegendItem{
					ID:    "id",
					SN:    "2",
					Start: "-80.75",
					End:   "-80.50",
					RGB:   "#220000",
				},
				IndicatorLegendItem{
					ID:    "id",
					SN:    "2",
					Start: "-80.50",
					End:   "-80.24",
					RGB:   "#330000",
				},
				IndicatorLegendItem{
					ID:    "id",
					SN:    "2",
					Start: "-80.25",
					End:   "",
					RGB:   "#440000",
				},
			},
		},
	}

	//获取文件内容
	points := tf.GetPoints(indicator)

	pointsJson, _ := json.Marshal(points)

	log.Println(string(pointsJson))
}


func _TestTestFileLock(t *testing.T) {

	//创建TestFilePool
	tfp := InitTestFilePool("", "3s", nil)
	if tfp.GetLock()==false {
		t.Error("GetLock failed")
		return
	}
	tfp.ReleaseLock()
	if tfp.GetLock()==false {
		t.Error("GetLock failed")
		return
	}
	time.Sleep(4 * time.Second)
	if tfp.GetLock()==false {
		t.Error("GetLock failed")
		return
	}
	time.Sleep(2 * time.Second)
	if tfp.GetLock()==true {
		t.Error("GetLock failed")
		return
	}
}