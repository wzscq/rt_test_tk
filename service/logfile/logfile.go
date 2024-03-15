package logfile

import (
	"path/filepath"
	"os"
	"rt_test_service/crv"
	"log"
	"errors"
	"rt_test_service/common"
)

const MODELID_LOG_FILE = "rt_log_file"

type LogFileItem struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	CreationTime string `json:"creationTime"`
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

func UpdateLogFilesToDB(files []LogFileItem,crvClient *crv.CRVClient,token string)(error){
	for _, file := range files {
		err:=UpdateLogFileToDB(file, crvClient,token)
		if err!=nil {
			return err
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