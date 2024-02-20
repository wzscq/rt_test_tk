package robot

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
)

var MODELID_TASK="rt_task"
var	queryTaskFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "name"},
	{"field": "robot_id"},
	{
		"field":"ue_tc",
		"fieldType":"one2many",
		"relatedModelID":"rt_task_ue_tc",
		"relatedField":"task_id",
		"fields":[]map[string]interface{}{
			{"field": "id"},
			{"field": "task_id"},
			{
				"field":"imei_id",
				"fieldType":"many2one",
				"relatedModelID":"rt_ue_imei",
				"fields":[]map[string]interface{}{
					{"field": "id"},
					{"field": "imei"},
					{"field": "imsi"},
				},
			},
			{
				"field": "tcs",
				"fieldType":"many2many",
				"relatedModelID":"rt_robot_test_case",
				"fields":[]map[string]interface{}{
					{"field": "id"},
					{"field": "tc_id"},
				},
			},
		},
	},
}

func GetTask(crvClient *crv.CRVClient,taskID string,token string,ftpConf *common.FtpConf)(map[string]interface{},int){
	filter:=map[string]interface{}{
		"id":taskID,
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_TASK,
		Fields:&queryTaskFields,
		Filter:&filter,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return nil,commonErr
	}

	if rsp.Error == true {
		log.Println("GetTask error:",rsp.ErrorCode,rsp.Message)
		return nil,common.ResultNoTask
	}

	if rsp.Result == nil {
		log.Println("GetTask error: no result in rsp.")
		return nil,common.ResultNoTask
	}

	resMap,ok:=rsp.Result.(map[string]interface{})
	if !ok {
		log.Println("GetRobotList error: can not convert item to map.")
		return nil,common.ResultNoTask
	}

	taskList,ok:=resMap["list"].([]interface{})
	if !ok || len(taskList) == 0 {
		log.Println("GetRobotList error: no task in rsp.")
		return nil,common.ResultNoTask
	}

	taskMap,ok:=taskList[0].(map[string]interface{})
	if !ok {
		log.Println("GetRobotList error: can not convert item to map.")
		return nil,common.ResultNoTask
	}

	return ConvertToSendTask(taskMap,ftpConf)
}

func ConvertToSendTask(task map[string]interface{},ftpConf *common.FtpConf)(map[string]interface{},int){
	sendTask:=map[string]interface{}{
		"testPlanId":task["id"],
		"exampleName":task["name"],
	}

	taskUE,ok:=task["ue_tc"].(map[string]interface{})
	if !ok {
		log.Println("ConvertToSendTask error: can not convert ue_tc to map.")
		return nil,common.ResultNoTaskUe
	}
	
	ueList:=taskUE["list"].([]interface{})
	if len(ueList) == 0 {
		log.Println("ConvertToSendTask error: no ue in task.")
		return nil,common.ResultNoTaskUe
	}

	details:=[]map[string]interface{}{
		map[string]interface{}{
		"tc_exe_id":task["id"],
		"tc_exe_name":task["name"],
		"ftp_conf":ftpConf,
		},
	}

	relation:=[]map[string]interface{}{}

	for _,ue:=range ueList {
		ueMap,_:=ue.(map[string]interface{})
		ueImei,_:=ueMap["imei_id"].(map[string]interface{})
		imeiList:=ueImei["list"].([]interface{})
		if len(imeiList) == 0 {
			return nil,common.ResultNoTaskUe
		}

		imeiMap,_:=imeiList[0].(map[string]interface{})

		ueDetail:=map[string]interface{}{
			"ue_identify":"imsi",
			"ue_id":imeiMap["imsi"],
		}

		ueTc,ok:=ueMap["tcs"].(map[string]interface{})
		if !ok {
			log.Println("ConvertToSendTask error: can not convert tcs to map.")
			return nil,common.ResultNoTaskUeTc
		}

		tcList:=ueTc["list"].([]interface{})
		if len(tcList) == 0 {
			log.Println("ConvertToSendTask error: no tc in ue.")
			return nil,common.ResultNoTaskUeTc
		}

		tcids:=[]string{}
		for _,tc:=range tcList {
			tcMap,_:=tc.(map[string]interface{})
			tcids=append(tcids,tcMap["tc_id"].(string))
		}
		ueDetail["tc_ids"]=tcids
		relation=append(relation,ueDetail)
	}
	details[0]["relation"]=relation
	sendTask["detail"]=details
	robotID:=task["robot_id"].(string)
	sendTask["robotId"]=robotID
	return sendTask,common.ResultSuccess
}