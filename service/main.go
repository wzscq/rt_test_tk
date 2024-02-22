package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"rt_test_service/common"
	//"rt_test_service/robot"
	"rt_test_service/mqtt"
	"rt_test_service/crv"
	"rt_test_service/device"
	"rt_test_service/testfile"
	"log"
	"time"
)

func main() {
	//设置log打印文件名和行号
  log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化配置文件
	conf:=common.InitConfig("conf/conf.json")

	//初始化时区
  var cstZone = time.FixedZone("CST", conf.Service.CSTZone*3600)
	time.Local = cstZone

	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	//crvClient 用于到crvframeserver的请求
	crvClient:=&crv.CRVClient{
		Server:conf.CRV.Server,
		Token:conf.CRV.Token,
		AppID:conf.CRV.AppID,
	}

	//初始化测试文件处理对象
	tfp:=testfile.InitTestFilePool(conf.TestFile.Path,conf.TestFile.IdleBeforeClose,crvClient)
	
	//初始化MQTT客户端
	mqttClient:=mqtt.MQTTClient{
		Broker:conf.Mqtt.Broker,
		User:conf.Mqtt.User,
		Password:conf.Mqtt.Password,
		UploadMeasurementMetrics:conf.Mqtt.UploadMeasurementMetrics,
		Handler:tfp,
		Port:conf.Mqtt.Port,
	}
	mqttClient.Init()

	/*rtc:=robot.RobotController{
		RobotClient:robotClient,
		CRVClient:crvClient,
		FtpConf:&conf.Ftp,
	}
	rtc.Bind(router)*/

	dc:=device.DeviceController{
		CRVClient:crvClient,
		MqttConf:&conf.Mqtt,
		FtpConf:&conf.Ftp,
		MQTTClient:&mqttClient,
	}
	dc.Bind(router)

	tc:=testfile.TestFileController{
		OutPath:conf.TestFile.Path,
	}
	tc.Bind(router)
	
	router.Run(conf.Service.Port)
}