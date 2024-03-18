package common

import (
	"log"
	"os"
	"encoding/json"
)

type OauthConf struct {
	URL string `json:"url"`
	AppKey string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

type RobotAPI struct {
	URL string `json:"url"`
}

type RobotClientConf struct {
	Oauth OauthConf `json:"oauth"`
	GetRobotList RobotAPI `json:"getRobotList"`
	GetCurrentRobotStatus RobotAPI `json:"getCurrentRobotStatus"`
	GetTestEquipmentStatus RobotAPI `json:"getTestEquipmentStatus"`
	SendTask RobotAPI `json:"sendTask"`
}

type serviceConf struct {
	Port string `json:"port"`
	CSTZone int `json:"cstZone"`
}

type RobotMQTTClientConf struct {
	Broker string `json:"broker"`
	User string `json:"user"`
	Password string `json:"password"`
}

type crvConf struct {
	Server string `json:"server"`
  AppID string `json:"appID"`
	Token string `json:"token"`
}

type DeviceClientConf struct {
	ServerUrl string `json:"serverUrl"`
}

type MqttConf struct {
	Broker string `json:"broker"`
	Port int `json:"port"`
	WSPort int `json:"wsPort"`
	Password string `json:"password"`
	User string `json:"user"`
	UploadMeasurementMetrics string `json:"uploadMeasurementMetrics"`
	UploadDeviceStatus string `json:"uploadDeviceStatus"`
	SendTestCaseTopic string `json:"sendTestCaseTopic"`
	DecodeResutlTopic string `json:"decodeResutlTopic"`
}

type TestfileConf struct {
	Path string `json:"path"`
	IdleBeforeClose string `json:"idleBeforeClose"`
}

type TestLogFileConf struct {
	Path string `json:"path"`
	DecoderUrl string `json:"decoderUrl"`
}

type FtpConf struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Folder string `json:"folder"`
}

type MapConf struct {
	Center []float64 `json:"center"`
  Zoom int `json:"zoom"`
  MaxZoom int `json:"maxZoom"` 
  MinZoom int  `json:"minZoom"` 
  URL string `json:"url"`
}

type Config struct {
	Service serviceConf `json:"service"`
	RobotMQTTClient RobotMQTTClientConf `json:"robotMQTTClient"`
	RobotClient RobotClientConf `json:"robotClient"`
	CRV crvConf `json:"crv"`
	Mqtt MqttConf `json:"mqtt"`
	Ftp FtpConf `json:"ftp"`
	TestFile TestfileConf `json:"testFile"`
	Map MapConf `json:"map"`
	DeviceClient DeviceClientConf `json:"deviceClient"`
	TestLogFile TestLogFileConf `json:"testLogFile"`
}

var gConfig Config

func InitConfig(confFile string)(*Config){
	log.Println("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := confFile
	filePtr, err := os.Open(fileName)
	if err != nil {
        log.Fatal("Open file failed [Err:%s]", err.Error())
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
	}
	log.Println("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}