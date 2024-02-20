package device

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
)

type ServerHeader struct {
	Token string  `json:"token"`
}

type DeviceReq struct {
	RobotID string `form:"robotId"`
}

var	queryTestCaseFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "name"},
	{"field": "version"},
	{"field": "rt_main_function_id"},
	{"field": "rt_function_method_id"},
	{
		"field": "params",
		"fieldType": "one2many",
		"relatedModelID": "rt_project_testcase_params",
		"relatedField": "rt_project_test_case_id",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "param_name"},
			{"field": "param_value"},
			{"field": "rt_project_test_case_id"},
		},
	},
}

var	queryRobotTestCaseFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
	{
		"field": "params",
		"fieldType": "one2many",
		"relatedModelID": "rt_robot_testcase_params",
		"relatedField": "rt_robot_test_case_id",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "version"},
			{"field": "rt_robot_test_case_id"},
		},
	},
}

var MODELID_PROJECT_TESTCASE="rt_project_test_case"
var MODELID_ROBOT_TESTCASE="rt_robot_test_case"

func GetCommitedTestCase(crvClient *crv.CRVClient,token string)(*common.CommonRsp){
	filter:=map[string]interface{}{
		"status":"1",
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_PROJECT_TESTCASE,
		Fields:&queryTestCaseFields,
		Filter:&filter,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return rsp
	}

	if rsp.Error == true {
		log.Println("GetCommitedTestCase error:",rsp.ErrorCode,rsp.Message)
		return rsp
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetCommitedTestCase error: no list in rsp.")
		rsp:=common.CreateResponse(common.CreateError(common.ResultNoCommitedTestCaseError,nil),nil)
		return rsp
	}

	if len(resLst)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultNoCommitedTestCaseError,nil),nil)
		return rsp
	}

	return rsp
}

func GetRobotTestCase(crvClient *crv.CRVClient,token string,robotId string)(*common.CommonRsp){
	filter:=map[string]interface{}{
		"rt_project_robot_id":robotId,
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT_TESTCASE,
		Fields:&queryRobotTestCaseFields,
		Filter:&filter,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		return rsp
	}

	return rsp
}

func UpdateRobotTestCase(crvClient *crv.CRVClient,token string,robotID string,tcList []interface{})(*common.CommonRsp){
	//获取并删除原有的机器人测试用例
	rsp:=GetRobotTestCase(crvClient,token,robotID)

	if rsp.Error==true {
		return rsp
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if ok && len(resLst)>0 {
		saveList:=[]map[string]interface{}{}
		for _,tcItem:=range resLst {
			tcItemMap,ok:=tcItem.(map[string]interface{})
			saveItemMap:=map[string]interface{}{
				"_save_type":"delete",
				"id":tcItemMap["id"],
				"version":tcItemMap["version"],
			}
			tcParams,ok:=tcItemMap["params"].(map[string]interface{})
			if ok {
				tcParamsList,ok:=tcParams["list"].([]interface{})
				if ok {
					saveParams:=map[string]interface{}{
						"fieldType":"one2many",
						"modelID": "rt_robot_testcase_params",
						"relatedField": "rt_robot_test_case_id",
					}
					saveParamsList:=[]map[string]interface{}{}
					for _,paramItem:=range tcParamsList {
						paramItemMap,_:=paramItem.(map[string]interface{})
						saveParamsItemMap:=map[string]interface{}{
							"_save_type":"delete",
							"id":paramItemMap["id"],
							"version":paramItemMap["version"],
						}
						saveParamsList=append(saveParamsList,saveParamsItemMap)
					}
					saveParams["list"]=saveParamsList
					saveItemMap["params"]=saveParams
				}
			}
			saveList=append(saveList,saveItemMap)
		}

		commonRep:=crv.CommonReq{
			ModelID:MODELID_ROBOT_TESTCASE,
			List:&saveList,
		}
	
		rsp,commonErr:=crvClient.Save(&commonRep,token)
		if commonErr!=common.ResultSuccess {
			rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
			return rsp
		}
	
		if rsp.Error == true {
			log.Println("UpdateRobotTestCase error:",rsp.ErrorCode,rsp.Message)
			return rsp
		}
	} else {
		log.Println("get robot test case error: no list in rsp.")
	}

	//添加新的测试用例到机器人测试用例
	saveList:=[]map[string]interface{}{}
	for _,tcItem:=range tcList {
		tcItemMap,ok:=tcItem.(map[string]interface{})
		saveItemMap:=map[string]interface{}{
			"_save_type":"create",
			"tc_id":tcItemMap["id"],
			"name":tcItemMap["name"],
			"rt_project_robot_id":robotID,
			"rt_main_function_id":tcItemMap["rt_main_function_id"],
			"rt_function_method_id":tcItemMap["rt_function_method_id"],
		}

		tcParams,ok:=tcItemMap["params"].(map[string]interface{})
		if ok {
			tcParamsList,ok:=tcParams["list"].([]interface{})
			if ok {
				saveParams:=map[string]interface{}{
					"fieldType":"one2many",
					"modelID": "rt_robot_testcase_params",
					"relatedField": "rt_robot_test_case_id",
				}
				saveParamsList:=[]interface{}{}
				for _,paramItem:=range tcParamsList {
					paramItemMap,_:=paramItem.(map[string]interface{})
					saveParamsItemMap:=map[string]interface{}{
						"_save_type":"create",
						"param_name":paramItemMap["param_name"],
						"param_value":paramItemMap["param_value"],
					}
					saveParamsList=append(saveParamsList,saveParamsItemMap)
				}
				saveParams["list"]=saveParamsList
				saveItemMap["params"]=saveParams
			}
		}
		saveList=append(saveList,saveItemMap)
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT_TESTCASE,
		List:&saveList,
	}

	rsp,commonErr:=crvClient.Save(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		return rsp
	}

	if rsp.Error == true {
		log.Println("UpdateRobotTestCase error:",rsp.ErrorCode,rsp.Message)
		return rsp
	}

	return nil
}

func GetTestCaseForDevice(tcList []interface{})(map[string]interface{}){
	list:=[]interface{}{}
	for _,tcItem:=range tcList {
		tcItemMap,ok:=tcItem.(map[string]interface{})
		dvtcItem:=map[string]interface{}{
			"tc_id":tcItemMap["id"],
			"tc_name":tcItemMap["name"],
			"main_function":tcItemMap["rt_main_function_id"],
			"method":tcItemMap["rt_function_method_id"],
		}
		if ok {
			tcParams,ok:=tcItemMap["params"].(map[string]interface{})
			if ok {
				tcParamsList,ok:=tcParams["list"].([]interface{})
				if ok {
					tc_params:=map[string]interface{}{}
					for _,paramItem:=range tcParamsList {
						paramItemMap,_:=paramItem.(map[string]interface{})
						tc_params[paramItemMap["param_name"].(string)]=paramItemMap["param_value"]
					}
					dvtcItem["tc_params"]=tc_params
				}
			}
		}
		list=append(list,dvtcItem)
	}

	res:=map[string]interface{}{
		"list":list,
	}
	return res
}