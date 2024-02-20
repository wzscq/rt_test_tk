package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"rt_test_service/common"
	"rt_test_service/robot"
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
	conf:=common.InitConfig()

	//初始化时区
    var cstZone = time.FixedZone("CST", conf.Service.CSTZone*3600)
	time.Local = cstZone

	//启动到机器人平台的mqtt连接
	/*robotMqttClient:=robot.RobotMQTTClient{
		Broker:conf.RobotMQTTClient.Broker,
		User:conf.RobotMQTTClient.User,
		Password:conf.RobotMQTTClient.Password,
	}
	robotMqttClient.Init()*/
	
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))

	//初始化机器人平台客户端
	robotClient:=&robot.RobotClient{
		Conf:&conf.RobotClient,
	}

	/*rsp,err:=robotClient.Oauth()
	if err==nil && rsp.Result!=nil {
		log.Println("RobotClient Oauth success with token:",rsp.Result.Token)
		rsp1,err:=robotClient.GetCurrentRobotStatus(rsp.Result.Token,"2bee174b7d7c36e9b98bd8772e66af5e")
		if err==nil && rsp1.Result!=nil {
			log.Println("RobotClient GetRobotList success with result:",rsp1)
		}
	}

	return;*/

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

	rtc:=robot.RobotController{
		RobotClient:robotClient,
		CRVClient:crvClient,
		FtpConf:&conf.Ftp,
	}
	rtc.Bind(router)

	dc:=device.DeviceController{
		CRVClient:crvClient,
		MqttConf:&conf.Mqtt,
		FtpConf:&conf.Ftp,
	}
	dc.Bind(router)

	tc:=testfile.TestFileController{
		OutPath:conf.TestFile.Path,
	}
	tc.Bind(router)
	
	router.Run(conf.Service.Port)
}