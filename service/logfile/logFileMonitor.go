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
	LogFileFromDecoderUrl string
}

func InitLogFileMonitor(conf *common.Config,crvClient *crv.CRVClient)(lfm *LogFileMonitor) {
	duration, _ := time.ParseDuration(conf.TestLogFile.ExpandTimeRange)
	lfm = &LogFileMonitor{
		LogFilePath: conf.TestLogFile.Path,
		ExpandTimeRange: duration,
		CRVClient: crvClient,
		LogFileFromDecoderUrl: conf.TestLogFile.LogFileFromDecoderUrl,
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
	files, err := GetLogFileListFromDecode(lfm.LogFileFromDecoderUrl)//GetLogFileList(lfm.LogFilePath)
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
	count:=GetTestLogByTime(file.CreationTime,lfm.CRVClient,"",lfm.ExpandTimeRange)
	if count == 0 {
		return false
	}

	//添加文件信息到数据库
	UpdateLogFileToDB(file, lfm.CRVClient, "")
	return true
}




