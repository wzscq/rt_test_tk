package device

import (
	"fmt"
	"rt_test_service/common"
	"rt_test_service/crv"
	"rt_test_service/mqtt"
	"testing"
	"time"
)

/*
{
	"details": {
		"exampleCode": "ftp_executor",
		"params": {
			"tc_params": {
				"method": "download",
				"user": "test",
				"passwd": "Happy@robot#1234",
				"remote_file": "/tmp/7gb",
				"server_url": "192.168.100.109",
				"mode": "bytime",
				"intervals": 10,
				"repeat": 1,
				"duration": 5,
				"no_data": 60,5
				"no_data_ping": 10,
				"thread": 3,
				"timeout": 20
			},
			"operator_info": {
				"band": "78",
				"ue_identify_type": "imsi",
				"ue_identify": "460011895631209",
				"freq": "1234",
				"netType": "LTE"
			},
			"tc_id": "tc001"
		}
	},
	"trigger": "start",
	"topic": "CommandResult"
}
*/

func TestGetTestCase(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}
	tc := GetTestCase("FTPUpload", "rt_test_tk_service", crvClient)
	if tc == nil {
		t.Error("GetTestCase error")
		return
	}
	fmt.Println(tc)
	if tc.ID != "FTPUpload" {
		t.Error("GetTestCase error")
	}
}

func TestSendTestCase(t *testing.T) {
	conf := common.InitConfig("../conf/conf.json")

	crvClient := &crv.CRVClient{
		Server: conf.CRV.Server,
		Token:  conf.CRV.Token,
		AppID:  conf.CRV.AppID,
	}

	tc := GetTestCase("001", conf.CRV.Token, crvClient)

	if tc == nil {
		t.Error("GetTestCase error")
		return
	}

	cmd:=GetTestCommand(tc)

	//初始化MQTT客户端
	mqttClient := mqtt.MQTTClient{
		Broker:                   conf.Mqtt.Broker,
		User:                     conf.Mqtt.User,
		Password:                 conf.Mqtt.Password,
		UploadMeasurementMetrics: conf.Mqtt.UploadMeasurementMetrics,
		Port:                     conf.Mqtt.Port,
	}
	mqttClient.Init()

	time.Sleep(2 * time.Second)

	err := SendTestCase(cmd, &mqttClient, conf.Mqtt.SendTestCaseTopic)
	if err != nil {
		t.Error("SendTestCase error")
	}
}
