package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	//"rt_test_service/robot"
	"log"
	"rt_test_service/crv"
	"rt_test_service/device"
	"rt_test_service/mqtt"
	"rt_test_service/testfile"
	"rt_test_service/logfile"
	"time"
)

func main() {
	//设置log打印文件名和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化配置文件
	conf := common.InitConfig("conf/conf.json")

	//初始化时区
	var cstZone = time.FixedZone("CST", conf.Service.CSTZone*3600)
	time.Local = cstZone

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	//crvClient 用于到crvframeserver的请求
	crvClient := &crv.CRVClient{
		Server: conf.CRV.Server,
		Token:  conf.CRV.Token,
		AppID:  conf.CRV.AppID,
	}

	//初始化测试文件处理对象
	//tfp := testfile.InitTestFilePool(conf.TestFile.Path, conf.TestFile.IdleBeforeClose, crvClient)

	//初始化MQTT客户端
	mqttClient := mqtt.MQTTClient{
		Broker:                   conf.Mqtt.Broker,
		User:                     conf.Mqtt.User,
		Password:                 conf.Mqtt.Password,
		UploadMeasurementMetrics: conf.Mqtt.UploadMeasurementMetrics,
		//Handler:                  tfp,
		Port:                     conf.Mqtt.Port,
	}
	mqttClient.Init()

	device.InitDeviceController(conf,crvClient,&mqttClient,router)
	logfile.InitLogFileController(conf,crvClient,router)
	
	tc := testfile.TestFileController{
		OutPath: conf.TestFile.Path,
	}
	tc.Bind(router)

	router.Static("/maptiles", "./maptiles")

	router.Run(conf.Service.Port)
}
