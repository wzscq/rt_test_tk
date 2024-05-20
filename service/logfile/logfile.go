package logfile

import (
	"path/filepath"
	"os"
	"rt_test_service/crv"
	"log"
	"errors"
	"rt_test_service/common"
	"time"
	"net/http"
	"encoding/json"
)

const MODELID_LOG_FILE = "rt_log_file"

type LogFileItem struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	CreationTime string `json:"creationTime"`
	Status string `json:"status"`
	DecodeID string `json:"decodeID"`
	DecodedFile string `json:"decodedFile"`
}

type DecoderFileItem struct {
	FileName string `json:"name"`
	CreationTime string `json:"creation_time"`
	ModificationTime string `json:"modification_time"`
}

// GetFileList returns a list of log files in the specified path.
// The filter parameter is used to filter the files by name.
func GetLogFileList(path string) ([]LogFileItem,error){
	var files []LogFileItem
	// Read all the files from the directory
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info!=nil && !info.IsDir() {
				
					fileinfo := LogFileItem {
						Name:         filepath.Base(info.Name()),
						Size:         info.Size(),
						CreationTime: info.ModTime().Format("2006-01-02 15:04:05"),
					}

					log.Println(fileinfo)

					files = append(files, fileinfo)
			}
			return nil
	})

	return files, err

}

func GetLogFileListFromDecode(url string)([]LogFileItem,error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("GetLogFileListFromDecode Get error", err)
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("GetLogFileListFromDecode Get Do request error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("GetLogFileListFromDecode Get StatusCode error", resp)
		return nil, errors.New("api returns a worng status code")
	}

	decoder := json.NewDecoder(resp.Body)
	fileList := []DecoderFileItem{}
	err = decoder.Decode(&fileList)
	if err != nil {
		log.Println("GetLogFileListFromDecode Get result decode failed error", err.Error())
		return nil, err
	}

	var files []LogFileItem
	//转换格式
	for _, file := range fileList {
		fileinfo := LogFileItem {
			Name:         file.FileName,
			CreationTime: GetCreateTime(file.FileName),
		}

		files = append(files, fileinfo)
	}
	

	log.Println("end DeviceClient Get success")
	return files, nil
}

func GetCreateTime(fileName string)(string){
	//文件名称格式为：20240519_073626_0014.qmdl2
	//转换为：2024-05-19 07:36:26
	return fileName[0:4]+"-"+fileName[4:6]+"-"+fileName[6:8]+" "+fileName[9:11]+":"+fileName[11:13]+":"+fileName[13:15]
}

func UpdateLogFilesToDB(files []LogFileItem,crvClient *crv.CRVClient,token string,expandTimeRange time.Duration)(error){
	for _, file := range files {
		//查询数据库中是否存在对应时间的测试日志
		count:=GetTestLogByTime(file.CreationTime,crvClient,token,expandTimeRange)
		if count > 0 {
			err:=UpdateLogFileToDB(file, crvClient,token)
			if err!=nil {
				return err
			}
		}
	}
	return nil
}

func UpdateLogFileToDB(file LogFileItem,crvClient *crv.CRVClient,token string)(error){
	//从数据库中查询这个文件
	fileInfo,err:=GetLogFileFromDB(file,crvClient,token)
	if err!=nil {
		return err
	}

	if fileInfo==nil {
		fileInfo=map[string]interface{}{
			"id":file.Name,
			"size":file.Size,
			"creation_time":file.CreationTime,
			"_save_type":"create",
		}
	} else {
		fileInfo["size"]=file.Size
		fileInfo["creation_time"]=file.CreationTime
		fileInfo["_save_type"]="update"
	}

	return SaveLogFileToDB(fileInfo,crvClient,token)
}

func SaveLogFileToDB(logFile map[string]interface{},crvClient *crv.CRVClient,token string)(error){
	commonRep := crv.CommonReq{
		ModelID: MODELID_LOG_FILE,
		List:    &[]map[string]interface{}{logFile},
	}

	_, commonErr := crvClient.Save(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return errors.New("SaveLogFileToDB error")
	}

	return nil
}

func GetLogFileFromDB(file LogFileItem,crvClient *crv.CRVClient,token string)(map[string]interface{},error){
	//从数据库中查询这个文件
	filter := map[string]interface{}{
		"id":map[string]interface{}{
			"Op.eq": file.Name,
		},
	}

	commonRep := crv.CommonReq{
		ModelID: MODELID_LOG_FILE,
		Fields:  &[]map[string]interface{}{
			{"field": "id"},
			{"field": "size"},
			{"field": "creation_time"},
			{"field": "version"},
		},
		Filter:  &filter,
	}

	rsp, commonErr := crvClient.Query(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return nil,errors.New("GetLogFileFromDB error")
	}

	if rsp.Error == true {
		log.Println("GetCommitedTestCase error:", rsp.ErrorCode, rsp.Message)
		return nil,errors.New("GetLogFileFromDB error")
	}

	resLst, ok := rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetCommitedTestCase error: no list in rsp.")
		return nil,nil
	}

	if len(resLst) == 0 {
		log.Println("GetCommitedTestCase error: no list in rsp.")
		return nil,nil
	}

	fileInfo,ok:=resLst[0].(map[string]interface{})
	if !ok {
		log.Println("GetCommitedTestCase error: item is not map.")
		return nil,nil
	}

	return fileInfo,nil
}

func DecodeLogFile(logFles *[]string,dc *DecoderClient,crvClient *crv.CRVClient,token string)(error){
	//获取解码器状态
	status,err:=dc.GetStatus()
	if err!=nil {
		return err
	}

	if status.Status!="ready" {
		return errors.New("当前解码器状态为："+status.Status+", 无法解析文件")
	}

	//解析文件
	res,err:=dc.DecodeFile(logFles)
	if err!=nil {
		return err
	}

	if res.Res!="accept" {
		return errors.New("解码器返回错误，错误信息："+res.Cause)
	}

	//创建解析记录
	return SaveDecodeRecordToDB(res,crvClient,token)
}

func DeleteAllLogFiles(crvClient *crv.CRVClient,token string) {
	//删除所有文件
	commonRep := crv.CommonReq{
		ModelID: MODELID_LOG_FILE,
		Filter: &map[string]interface{}{
			"create_time": map[string]interface{}{
				"Op.lt": time.Now().Format("2000-01-01 00:00:00"),
			},
		},
		SelectedRowKeys: &[]string{},
		SelectAll:true,
	}

	crvClient.Delete(&commonRep, token)
}

func DeleteLogFileByName(name string,crvClient *crv.CRVClient,token string) {
	//删除所有文件
	commonRep := crv.CommonReq{
		ModelID: MODELID_LOG_FILE,
		SelectedRowKeys: &[]string{name},
	}

	crvClient.Delete(&commonRep, token)
}

func GetTestLogByTime(timeStr string,crvClient *crv.CRVClient,token string,expandTimeRange time.Duration)(int) {
	//string to time
	startTime, _:= time.Parse("2006-01-02 15:04:05", timeStr)
	startTime=startTime.Add(time.Duration(expandTimeRange))

	endTime,_:= time.Parse("2006-01-02 15:04:05", timeStr)
	endTime=endTime.Add(-time.Duration(expandTimeRange))

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

	rsp, commonErr := crvClient.Query(&commonRep, token)
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