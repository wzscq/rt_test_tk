package testfile

import (
	"log"
	"sync"
	"time"
	"rt_test_service/crv"
	"encoding/json"
	"fmt"
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
	OutPath string
	Pool map[string]*TestFile
	IdleBeforeClose time.Duration
	Mutex sync.Mutex
	CRVClient *crv.CRVClient
}

func InitTestFilePool(outPath string,idleBeforeClose string,crvClient *crv.CRVClient)(*TestFilePool){
	duration, _ := time.ParseDuration(idleBeforeClose)

	tfp:=&TestFilePool{
		OutPath:outPath,
		Pool:make(map[string]*TestFile),
		IdleBeforeClose:duration,
		CRVClient:crvClient,
	}

	//启动扫描线程
	go tfp.Scan()

	return tfp
}

func (tfp *TestFilePool)createCacheRecord(tf *TestFile){
	//将unix时间戳转换为字符串
	startTime:=time.Unix(tf.StartTime,0).Format("2006-01-02 15:04:05")

	testFileMap:=map[string]interface{}{
		"device_id":tf.DeviceID,
		"timestamp":tf.TimeStamp,
		"start_time":startTime,
		"line_count":tf.LineCount,
		"_save_type":"create",
	}

	commonRep:=crv.CommonReq{
		ModelID:"rt_cache_test_file",
		List:&[]map[string]interface{}{testFileMap},
	}

	tfp.CRVClient.Save(&commonRep,"")
}

func (tfp *TestFilePool)Scan(){
	//间隔IdleBeforeClose秒扫描一次，对于没有写入的文件，关闭
	for {
		time.Sleep(time.Duration(tfp.IdleBeforeClose))
		log.Println("TestFilePool.Scan ...")
		tfp.Mutex.Lock()
		for _,tf:=range tfp.Pool {
			if tf.LineCount==tf.lastLineCount {
				tf.Close()
				log.Println("TestFilePool.Scan close test file with deviceID:"+tf.DeviceID)
				delete(tfp.Pool,tf.DeviceID)
				tfp.createCacheRecord(tf)
			} else {
				tf.lastLineCount=tf.LineCount
			}
		}
		tfp.Mutex.Unlock()
	}
}

func (tfp *TestFilePool) WriteDeviceTestLine(deviceID,line string){
	//这里需要枷锁做并发控制
	tfp.Mutex.Lock()
	defer tfp.Mutex.Unlock()

	tf:=tfp.Pool[deviceID]
	if tf == nil {
		tf=tfp.CreateTestFile(deviceID,line)
		tfp.Pool[deviceID]=tf
	}
	
	if tf == nil {
		return 
	}

	tf.WriteLine(line)
}

func GetLineTimeStamp(line string)(string){
	//解析line，获取时间戳
	LineData:=LineData{}
	err:=json.Unmarshal([]byte(line),&LineData)
	if err != nil {
		log.Println("GetLineTimeStamp Unmarshal failed:",err)
		return ""
	}

	if len(LineData.Data)==0 {
		log.Println("GetLineTimeStamp LineData.Data is empty")
		return ""
	}

	return fmt.Sprintf("%d",LineData.Data[0].CaseProgress.SessionID)
}

func (tfp *TestFilePool)CreateTestFile(deviceID,line string)(*TestFile){
	timeStamp:=GetLineTimeStamp(line)
	return GetTestFile(tfp.OutPath,deviceID,timeStamp)
}

func (tfp *TestFilePool)DealDeviceTestMessage(deviceID,line string){
	tfp.WriteDeviceTestLine(deviceID,line)
}