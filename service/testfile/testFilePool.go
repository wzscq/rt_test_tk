package testfile

import (
	"encoding/json"
	"log"
	"rt_test_service/crv"
	"rt_test_service/device"
	"sync"
	"time"
	"rt_test_service/common"
)

type CaseProgress struct {
	SessionID int64 `json:"session_id"`
}

type LineDataItem struct {
	CaseProgress CaseProgress `json:"case_progress"`
}

type LineData struct {
	Data []LineDataItem `json:"data"`
}

type TestFilePool struct {
	OutPath         string
	Pool            map[string]*TestFile
	IdleBeforeClose time.Duration
	Mutex           sync.Mutex
	CRVClient       *crv.CRVClient
	IsRunning       bool 
	StartTime       time.Time
	CmdSender       device.CmdSender
	LastGPSRec	    *GPSRec
}

type TestData struct {
	GPS map[string]interface{} `json:"gps"`
	Result interface{} `json:"result"`
}

type ReportData struct {
	ExampleCode string `json:"exampleCode"`
    MsgType string `json:"msg_type"`
	TestData TestData `json:"testData"`
}

func InitTestFilePool(outPath string, idleBeforeClose string, crvClient *crv.CRVClient) *TestFilePool {
	duration, _ := time.ParseDuration(idleBeforeClose)

	tfp := &TestFilePool{
		OutPath:         outPath,
		Pool:            make(map[string]*TestFile),
		IdleBeforeClose: duration,
		CRVClient:       crvClient,
	}

	//启动扫描线程
	if crvClient != nil {
		go tfp.Scan()
	}

	return tfp
}

func (tfp *TestFilePool) createCacheRecord(tf *TestFile) {
	//将unix时间戳转换为字符串
	startTime := time.Unix(tf.StartTime, 0).Format("2006-01-02 15:04:05")

	testFileMap := map[string]interface{}{
		"device_id":  tf.DeviceID,
		"timestamp":  tf.TimeStamp,
		"start_time": startTime,
		"line_count": tf.LineCount,
		"_save_type": "create",
	}

	commonRep := crv.CommonReq{
		ModelID: "rt_cache_test_file",
		List:    &[]map[string]interface{}{testFileMap},
	}

	tfp.CRVClient.Save(&commonRep, "")
}

func (tfp *TestFilePool) Scan() {
	//间隔IdleBeforeClose秒扫描一次，对于没有写入的文件，关闭
	for {
		time.Sleep(time.Duration(tfp.IdleBeforeClose))
		log.Println("TestFilePool.Scan ...")
		tfp.Mutex.Lock()
		for _, tf := range tfp.Pool {
			if tf.LineCount == tf.lastLineCount {
				if tf.sameLineCount>0 {
					//再一定时间内没有收到新的数据，认为测试已经停止，需要关闭测试记录文件
					tf.Close()
					log.Println("TestFilePool.Scan close test file with deviceID:" + tf.DeviceID)
					delete(tfp.Pool, tf.DeviceID)
					tfp.createCacheRecord(tf)
					//这里检查一下测试如果是循环测试，则继续下发下一个测试指令
					if(tfp.CmdSender!=nil){
						log.Println("TestFilePool.Scan send next cmd")
						tfp.CmdSender.SendCmd()
					} else {
						//释放锁,认为测试已经结束
						tfp.ReleaseLock()
					}
				} else {
					tf.sameLineCount++
				}
			} else {
				tf.lastLineCount = tf.LineCount
				tf.sameLineCount=0
			}
		}
		tfp.Mutex.Unlock()
	}
}

func (tfp *TestFilePool) WriteDeviceTestLine(deviceID, line string) {
	//这里需要枷锁做并发控制
	tfp.Mutex.Lock()
	defer tfp.Mutex.Unlock()

	tf := tfp.Pool[deviceID]
	if tf == nil {
		tf = tfp.CreateTestFile(deviceID)
		tfp.Pool[deviceID] = tf
	}

	if tf == nil {
		return
	}

	tf.WriteLine(line)
}

func GetLineTimeStamp() string {
	//获取当前时间戳
	return time.Now().Format("20060102150405")
}

func (tfp *TestFilePool) CreateTestFile(deviceID string) *TestFile {
	timeStamp := GetLineTimeStamp()
	return GetTestFile(tfp.OutPath, deviceID, timeStamp)
}

func (tfp *TestFilePool) SaveResult(deviceID string) {
	tfp.Mutex.Lock()
	defer tfp.Mutex.Unlock()

	tf := tfp.Pool[deviceID]
	if tf != nil {
		tf.Close()
		log.Println("SaveResult close test file with deviceID:" + tf.DeviceID)
		delete(tfp.Pool, tf.DeviceID)
		tfp.createCacheRecord(tf)
	}
}

func (tfp *TestFilePool) AddResult(deviceID string, result map[string]interface{}) {
	//这里需要枷锁做并发控制
	resultByte, err := json.MarshalIndent(result, "", "    ")
	var resultString string
	if err != nil {
		resultString=""
	} else {
		resultString = string(resultByte)
	}

	tfp.Mutex.Lock()
	defer tfp.Mutex.Unlock()

	tf := tfp.Pool[deviceID]
	if tf != nil {
		tf.AddResult(resultString)
		log.Println("AddResult to test file with deviceID:" + tf.DeviceID)
		//delete(tfp.Pool, tf.DeviceID)
		//tfp.createCacheRecord(tf)
	}
}

func (tfp *TestFilePool) AddProcessResult(deviceID string,result map[string]interface{}) {
	//测试结果中补充GPS信息
	if tfp.LastGPSRec!=nil {
		result["gps"]=tfp.LastGPSRec
	}
	
	//这里需要枷锁做并发控制
	resultByte, err := json.MarshalIndent(result, "", "    ")
	var resultString string
	if err != nil {
		resultString=""
	} else {
		resultString = string(resultByte)
	}

	tfp.Mutex.Lock()
	defer tfp.Mutex.Unlock()

	tf := tfp.Pool[deviceID]
	if tf == nil {
		tf = tfp.CreateTestFile(deviceID)
		tfp.Pool[deviceID] = tf
	}

	if tf == nil {
		return
	}

	//获取第一个tf
	//for _, tf := range tfp.Pool {
	//	if tf != nil {
			tf.AddResult(resultString)
			//log.Println("AddResult to test file with deviceID:" + tf.DeviceID)
			//delete(tfp.Pool, tf.DeviceID)
			//tfp.createCacheRecord(tf)
	//	}
	//	break
	//}
}

func (tfp *TestFilePool) GetLock()(bool) {
	if tfp.IsRunning == true {
		//执行任务在规定时间内未更新时间则可能任务意外中断，需要重置
		if time.Now().Sub(tfp.StartTime) < tfp.IdleBeforeClose {
			return false
		}
	}
	tfp.IsRunning = true
	tfp.StartTime = time.Now()
	return true
}

func (tfp *TestFilePool) ReleaseLock() {
	tfp.IsRunning = false
}

func (tfp *TestFilePool) SetCmdSender(cmdSender device.CmdSender){
	tfp.CmdSender=cmdSender
}

func (tfp *TestFilePool) SaveGPS(gpsJsonStr string){
	SaveGPS(gpsJsonStr,tfp.CRVClient)
}

func (tfp *TestFilePool) HandleReportResult(report string) {
	//收到消息时更新时间，防止任务意外中断
	tfp.StartTime=time.Now()

	//decode to reportData
	reportData := ReportData{}
	err := json.Unmarshal([]byte(report), &reportData)
	if err != nil {
		log.Println("HandleReportResult Unmarshal failed:", err)
		return
	}

	if reportData.TestData.GPS!=nil {
		line, _ := json.Marshal(reportData.TestData.GPS)
		tfp.WriteDeviceTestLine(reportData.ExampleCode, string(line))
		tfp.LastGPSRec=GetGPSRec(string(line))
		tfp.SaveGPS(string(line))
	}

	//如果msg_type是finaly，则结束测试，保存文件
	if reportData.MsgType == "finally" {
		resultMap,ok:=reportData.TestData.Result.(map[string]interface{})
		if ok {
			tfp.AddResult(reportData.ExampleCode, resultMap)
		} else {
			resultString,ok:=reportData.TestData.Result.(string)
			if ok {
				resultMap:=map[string]interface{}{"result":resultString}
				tfp.AddResult(reportData.ExampleCode, resultMap)
			}
		}

		//这里检查一下测试如果是循环测试，则继续下发下一个测试指令
		if(tfp.CmdSender!=nil){
			log.Println("TestFilePool.Scan send next cmd")
			tfp.CmdSender.SendCmd()
		} else {
			//释放锁,认为测试已经结束
			tfp.SaveResult(reportData.ExampleCode)
			tfp.ReleaseLock()
		}
	}
}

func (tfp *TestFilePool) HandleResult(report string) {
	resultMap:=map[string]interface{}{}
	err:=json.Unmarshal([]byte(report), &resultMap)
	if err!=nil {
		log.Println("HandleResult Unmarshal failed:", err)
		return
	}

	deviceID,_:=resultMap["exampleCode"].(string)
	
	tfp.AddProcessResult(deviceID,resultMap)
}


/*func (tfp *TestFilePool) HandleReportResult(report string) {
	//收到消息时更新时间，防止任务意外中断
	tfp.StartTime=time.Now()

	//decode to reportData
	reportData := ReportData{}
	err := json.Unmarshal([]byte(report), &reportData)
	if err != nil {
		log.Println("HandleReportResult Unmarshal failed:", err)
		return
	}

	if reportData.TestData.GPS!=nil {
		line, _ := json.Marshal(reportData.TestData.GPS)
		tfp.WriteDeviceTestLine(reportData.ExampleCode, string(line))
	}

	//如果msg_type是finaly，则结束测试，保存文件
	if reportData.MsgType == "finally" {
		resultMap,ok:=reportData.TestData.Result.(map[string]interface{})
		if ok {
			tfp.SaveResult(reportData.ExampleCode, resultMap)
		} else {
			resultString,ok:=reportData.TestData.Result.(string)
			if ok {
				resultMap:=map[string]interface{}{"result":resultString}
				tfp.SaveResult(reportData.ExampleCode, resultMap)
			}
		}

		//这里检查一下测试如果是循环测试，则继续下发下一个测试指令
		if(tfp.CmdSender!=nil){
			log.Println("TestFilePool.Scan send next cmd")
			tfp.CmdSender.SendCmd()
		} else {
			//释放锁,认为测试已经结束
			tfp.ReleaseLock()
		}
	}	
}*/

func GetTestFileFromDB(id,token string,crvClient *crv.CRVClient)(*TestFile,int) {
	commonRep := crv.CommonReq{
		ModelID: "rt_cache_test_file",
		Filter: &map[string]interface{}{
			"id": map[string]interface{}{
				"Op.eq": id,
			},
		},
		Fields: &[]map[string]interface{}{
			{"field": "device_id"},
			{"field": "timestamp"},
		},
	}

	rsp,errorCode:=crvClient.Query(&commonRep, token)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	if rsp.Error {
		log.Println("GetTestFileFromDB error:",rsp.ErrorCode,rsp.Message)
		return nil,common.ResultDownloadFileError
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetTestFileFromDB error: no list in rsp.")
		return nil,common.ResultDownloadFileError
	}

	if len(resLst)==0 {
		log.Println("GetTestFileFromDB error: no list in rsp.")
		return nil,common.ResultDownloadFileError
	}

	recInfo,ok:=resLst[0].(map[string]interface{})
	if !ok {
		log.Println("GetTestFileFromDB error: item is not map.")
		return nil,common.ResultDownloadFileError
	}

	deviceID,_:=recInfo["device_id"].(string)
	timeStamp,_:=recInfo["timestamp"].(string)

	tf:=&TestFile{
		DeviceID:    deviceID,
		TimeStamp:   timeStamp,
	}

	return tf,common.ResultSuccess
}
