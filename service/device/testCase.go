package device

import (
	"log"
	"rt_test_service/crv"
	"rt_test_service/common"
	"rt_test_service/mqtt"
	"encoding/json"
	"errors"
)

type TestCase struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Function string `json:"function"`
	Method string `json:"method"`
	Params map[string]interface{} `json:"params"`
}

var	testCaseFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "name"},
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

func GetTestCase(id string,token string,crvClient *crv.CRVClient)(*TestCase){
	log.Println("GetTestCase")

	filter:=map[string]interface{}{
		"id":id,
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_PROJECT_TESTCASE,
		Fields:&testCaseFields,
		Filter:&filter,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return nil
	}

	if rsp.Error == true {
		log.Println("GetTestCase error:",rsp.ErrorCode,rsp.Message)
		return nil
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetTestCase error: no list in rsp.")
		return nil
	}

	if len(resLst)==0 {
		log.Println("GetTestCase error: list in rsp is empty.")
		return nil
	}

	tc:=TestCase{}
	tcMap:=resLst[0].(map[string]interface{})
	tc.ID,_=tcMap["id"].(string)
	tc.Name,_=tcMap["name"].(string)
	tc.Function,_=tcMap["rt_main_function_id"].(string)
	tc.Method,_=tcMap["rt_function_method_id"].(string)

	params,_:=tcMap["params"].([]interface{})
	for _,param:=range params {
		paramMap:=param.(map[string]interface{})
		paramName,_:=paramMap["param_name"].(string)
		paramValue,_:=paramMap["param_value"].(string)
		tc.Params[paramName]=paramValue
	}

	return &tc
}

func SendTestCase(tc *TestCase,mqttClient *mqtt.MQTTClient,topic string)(error) {
	log.Println("SendTestCase")
	//convert tc to json
	tcJson,err:=json.Marshal(tc)
	if err!=nil {
		log.Println("SendTestCase error: convert tc to json failed")
		log.Println(err)
		return err
	}

	//send tc to mqtt
	errCode:=mqttClient.Publish(topic,string(tcJson))
	if errCode!=common.ResultSuccess {
		log.Println("SendTestCase error: mqtt publish failed")
		//create a error 
		return errors.New("mqtt publish failed")
	}

	return nil
}