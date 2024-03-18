package logfile

import (
	"rt_test_service/crv"
	"log"
	"errors"
	"rt_test_service/common"
)

const MODELID_DECODE_REC = "rt_decode_record"

func GetDeodeRecordFromDB(decodeID int64,crvClient *crv.CRVClient,token string)(map[string]interface{},error){
	//从数据库中查询这个文件
	filter := map[string]interface{}{
		"id":map[string]interface{}{
			"Op.eq": decodeID,
		},
	}

	commonRep := crv.CommonReq{
		ModelID: MODELID_DECODE_REC,
		Fields:  &[]map[string]interface{}{
			{"field": "id"},
			{"field": "version"},
		},
		Filter:  &filter,
	}

	rsp, commonErr := crvClient.Query(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return nil,errors.New("GetDeodeRecordFromDB error")
	}

	if rsp.Error == true {
		log.Println("GetDeodeRecordFromDB error:", rsp.ErrorCode, rsp.Message)
		return nil,errors.New("GetDeodeRecordFromDB error")
	}

	resLst, ok := rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetDeodeRecordFromDB error: no list in rsp.")
		return nil,nil
	}

	if len(resLst) == 0 {
		log.Println("GetDeodeRecordFromDB error: no list in rsp.")
		return nil,nil
	}

	recInfo,ok:=resLst[0].(map[string]interface{})
	if !ok {
		log.Println("GetDeodeRecordFromDB error: item is not map.")
		return nil,nil
	}

	return recInfo,nil
}

func SaveDecodeRecordToDB(decodeRes *DecodeFileResponse,crvClient *crv.CRVClient,token string)(error){
	commonRep := crv.CommonReq{
		ModelID: MODELID_DECODE_REC,
		List:    &[]map[string]interface{}{
			map[string]interface{}{
				"id":decodeRes.Unixtime,
				"_save_type":"create",
			},
		},
	}

	_, commonErr := crvClient.Save(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return errors.New("保存解码记录信息到数据库是发生错误")
	}

	return nil
}

func UpdateDecodeStatus(decodeStatus *DecodeStatus,crvClient *crv.CRVClient,token string)(error){
	rec,err:=GetDeodeRecordFromDB(decodeStatus.DecodeID,crvClient,token)
	if err!=nil{
		return err
	}

	if rec==nil{
		return errors.New("UpdateDecodeStatus error: no record in db")
	}
	
	commonRep := crv.CommonReq{
		ModelID: MODELID_DECODE_REC,
		List:    &[]map[string]interface{}{
			map[string]interface{}{
				"id":rec["id"],
				"version":rec["version"],
				"status":decodeStatus.Status,
				"phase":decodeStatus.CurrentPhase,
				"decoded_file":decodeStatus.DecodeFiles,
				"_save_type":"update",
			},
		},
	}

	_, commonErr := crvClient.Save(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return errors.New("保存解码记录信息到数据库是发生错误")
	}

	return nil
}