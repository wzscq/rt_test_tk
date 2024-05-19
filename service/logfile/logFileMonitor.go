package logfile

import (
	"rt_test_service/crv"
	"github.com/wzscq/taskschedule"
	"time"
	"log"
	"rt_test_service/common"
)

type LogFileMonitor struct {
	CRVClient *crv.CRVClient
	LogFilePath string
	ExpandTimeRange time.Duration
	DealedFiles map[string]bool
}

func InitLogFileMonitor(conf *common.Config,crvClient *crv.CRVClient)(lfm *LogFileMonitor) {
	duration, _ := time.ParseDuration(conf.TestLogFile.ExpandTimeRange)
	lfm = &LogFileMonitor{
		LogFilePath: conf.TestLogFile.Path,
		ExpandTimeRange: duration,
		CRVClient: crvClient,
	}

	//清空数据库
	DeleteAllLogFiles(crvClient,"")

	//启动扫描任务
	schedule := &taskschedule.Schedule{
		Duration: conf.TestLogFile.ScanDuration,
		RunTime: time.Now().Format("15:04:05"),
	}

	taskschedule.RunTask(schedule,lfm)

	return lfm
}

func (lfm *LogFileMonitor) Run() {
	//读取文件列表
	files, err := GetLogFileList(lfm.LogFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	dealedFiles := map[string]bool{}
	//循环处理每个文件
	for _, file := range files {
		//更新文件信息到数据库
		isDeal:=lfm.UpdateLogFile(file)
		if isDeal == true {
			dealedFiles[file.Name] = true
		}
	}

	//删除过期的文件
	for name, _ := range lfm.DealedFiles {
		if !dealedFiles[name] {
			DeleteLogFileByName(name, lfm.CRVClient, "")
		}
	}

	lfm.DealedFiles = dealedFiles
}

func (lfm *LogFileMonitor) UpdateLogFile(file LogFileItem)(bool) {
	//判断文件是否已经检查过
	if lfm.DealedFiles[file.Name] {
		return true
	}

	//查询数据库中是否存在对应时间的测试日志
	count:=lfm.GetTestLogByTime(file.CreationTime)
	if count == 0 {
		return false
	}

	//添加文件信息到数据库
	UpdateLogFileToDB(file, lfm.CRVClient, "")
	return true
}

func (lfm *LogFileMonitor) GetTestLogByTime(timeStr string)(int) {
	//string to time
	startTime, _:= time.Parse("2006-01-02 15:04:05", timeStr)
	startTime=startTime.Add(time.Duration(lfm.ExpandTimeRange))

	endTime,_:= time.Parse("2006-01-02 15:04:05", timeStr)
	endTime=endTime.Add(-time.Duration(lfm.ExpandTimeRange))

	log.Println("GetTestLogByTime:", timeStr,startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))


	filter := map[string]interface{}{
		"start_time":map[string]interface{}{
			"Op.lte": startTime.Format("2006-01-02 15:04:05"),
		},
		"update_time":map[string]interface{}{
			"Op.gte": endTime.Format("2006-01-02 15:04:05"),
		},
	}

	commonRep := crv.CommonReq{
		ModelID: "rt_cache_test_file",
		Fields:  &[]map[string]interface{}{
			{"field": "id"},
		},
		Filter:  &filter,
	}

	rsp, commonErr := lfm.CRVClient.Query(&commonRep, "")
	if commonErr != common.ResultSuccess {
		return 0
	}

	if rsp.Error == true {
		log.Println("GetCommitedTestCase error:", rsp.ErrorCode, rsp.Message)
		return 0
	}

	resLst, ok := rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetTestLogByTime error: no list in rsp.")
		return 0
	}

	return len(resLst)
}



