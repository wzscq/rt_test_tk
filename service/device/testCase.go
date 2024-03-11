package device

import (
	"encoding/json"
	"errors"
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
	"rt_test_service/mqtt"
)

type TestCommandParams struct {
	Params map[string]interface{} `json:"tc_params"`
	OperatorInfo map[string]interface{} `json:"operator_info"`
	TCID string `json:"tc_id"`
}

type TestCommandDetail struct {
	ExampleCode string `json:"exampleCode"`
	Params TestCommandParams `json:"params"`
}

type TestCommand struct {
	Details TestCommandDetail `json:"details"`
	Trigger string `json:"trigger"`
	Topic string `json:"topic"`
}

type TestCase struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Function string                 `json:"function"`
	Method   string                 `json:"method"`
	Params   map[string]interface{} `json:"params"`
}

var testCaseFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "name"},
	{"field": "rt_main_function_id"},
	{"field": "rt_function_method_id"},
	{
		"field":          "params",
		"fieldType":      "one2many",
		"relatedModelID": "rt_project_testcase_params",
		"relatedField":   "rt_project_test_case_id",
		"fields": []map[string]interface{}{
			{"field": "id"},
			{"field": "param_name"},
			{"field": "param_value"},
			{"field": "rt_project_test_case_id"},
		},
	},
}

func GetTestCase(id string, token string, crvClient *crv.CRVClient) *TestCase {
	log.Println("GetTestCase")

	filter := map[string]interface{}{
		"id": id,
	}

	commonRep := crv.CommonReq{
		ModelID: MODELID_PROJECT_TESTCASE,
		Fields:  &testCaseFields,
		Filter:  &filter,
	}

	rsp, commonErr := crvClient.Query(&commonRep, token)
	if commonErr != common.ResultSuccess {
		return nil
	}

	if rsp.Error == true {
		log.Println("GetTestCase error:", rsp.ErrorCode, rsp.Message)
		return nil
	}

	resLst, ok := rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetTestCase error: no list in rsp.")
		return nil
	}

	if len(resLst) == 0 {
		log.Println("GetTestCase error: list in rsp is empty.")
		return nil
	}

	tc := TestCase{}
	tcMap := resLst[0].(map[string]interface{})
	tc.ID, _ = tcMap["id"].(string)
	tc.Name, _ = tcMap["name"].(string)
	tc.Function, _ = tcMap["rt_main_function_id"].(string)
	tc.Method, _ = tcMap["rt_function_method_id"].(string)
	tc.Params = make(map[string]interface{})
	paramMap, _ := tcMap["params"].(map[string]interface{})
	params, _ := paramMap["list"].([]interface{})
	for _, param := range params {
		paramMap := param.(map[string]interface{})
		paramName, _ := paramMap["param_name"].(string)
		paramValue, _ := paramMap["param_value"].(string)
		log.Println("paramName:", paramName, "paramValue:", paramValue)
		tc.Params[paramName] = paramValue
	}

	return &tc
}

func GetTestCommand(testCase *TestCase) *TestCommand {
	testCommandParams:=TestCommandParams{
		Params:testCase.Params,
		OperatorInfo:map[string]interface{}{
			"band": "78",
			"ue_identify_type": "imsi",
			"ue_identify": "460011895631209",
			"freq": "1234",
			"netType": "LTE",
		},
		TCID:testCase.ID,
	}

	testCommandDetail:=TestCommandDetail{
		ExampleCode:testCase.Name,
		Params:testCommandParams,
	}

	testCommand:=TestCommand{
		Trigger:"start",
		Topic:"CommandResult",
		Details:testCommandDetail,
	}

	return &testCommand
}

func SendTestCase(cmd *TestCommand, mqttClient *mqtt.MQTTClient, topic string) error {
	log.Println("SendTestCase")
	//convert tc to json
	tcJson, err := json.Marshal(cmd)
	if err != nil {
		log.Println("SendTestCase error: convert tc to json failed")
		log.Println(err)
		return err
	}

	//send tc to mqtt
	errCode := mqttClient.Publish(topic, string(tcJson))
	if errCode != common.ResultSuccess {
		log.Println("SendTestCase error: mqtt publish failed")
		//create a error
		return errors.New("mqtt publish failed")
	}

	return nil
}
