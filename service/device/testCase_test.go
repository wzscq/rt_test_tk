package device

import (
	"fmt"
	"testing"
	"rt_test_service/crv"
	"rt_test_service/common"
	"rt_test_service/mqtt"
	"time"
)

func TestGetTestCase(t *testing.T) {
	crvClient:=&crv.CRVClient{
		Server:"http://localhost:8200",
		Token:"rt_test_tk_service",
		AppID:"",
	}
	tc:=GetTestCase("FTPUpload","rt_test_tk_service",crvClient)
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
	conf:=common.InitConfig("../conf/conf.json")

	crvClient:=&crv.CRVClient{
		Server:conf.CRV.Server,
		Token:conf.CRV.Token,
		AppID:conf.CRV.AppID,
	}

	tc:=GetTestCase("FTPUpload",conf.CRV.Token,crvClient)

	if tc == nil {
		t.Error("GetTestCase error")
		return
	}

	//初始化MQTT客户端
	mqttClient:=mqtt.MQTTClient{
		Broker:conf.Mqtt.Broker,
		User:conf.Mqtt.User,
		Password:conf.Mqtt.Password,
		UploadMeasurementMetrics:conf.Mqtt.UploadMeasurementMetrics,
		Port:conf.Mqtt.Port,
	}
	mqttClient.Init()

	time.Sleep(2 * time.Second)

	err:=SendTestCase(tc,&mqttClient,conf.Mqtt.SendTestCaseTopic)
	if err != nil {
		t.Error("SendTestCase error")
	}
}