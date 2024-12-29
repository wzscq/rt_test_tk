package mqtt

import (
	"log"
	"time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"rt_test_service/common"
	"strconv"
)

type EventHandler interface {
	HandleDecodeResult(result string)
}

type ReportHandler interface {
	HandleReportResult(result string)
	HandleResult(result string)
}

type MQTTClient struct {
	Broker string
	User string
	Port int
	Password string
	UploadMeasurementMetrics string
	Handler EventHandler
	Client mqtt.Client
	DecodeResutlTopic string
	ReportHandler ReportHandler
}

const (
	PINT_RESULT = "ping_result"
	ATTACH_RESULT = "attach_result"
	TCP_RESULT = "tcp_result"
)

func (mqc *MQTTClient) getClient()(mqtt.Client){
	timeStamp:=time.Now().Unix()
	clientID:="rt_test_service_"+strconv.FormatInt(timeStamp,10)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://"+mqc.Broker+":"+strconv.Itoa(mqc.Port))
	opts.SetClientID(clientID)
	opts.SetUsername(mqc.User)
	opts.SetPassword(mqc.Password)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(mqc.messagePublishHandler)
	opts.OnConnect = mqc.connectHandler
	opts.OnConnectionLost = mqc.connectLostHandler
	opts.OnReconnecting = mqc.reconnectingHandler

	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("can not connect to mqtt server ",token.Error)
		return nil
	}
	return client
}

func (mqc *MQTTClient) connectHandler(client mqtt.Client){
	log.Println("MQTTClient connectHandler connect status: ",client.IsConnected())
	mqc.Client=client
	if client.IsConnected() {
		log.Println("MQTTClient Subscribe topic:"+mqc.UploadMeasurementMetrics)
		client.Subscribe(mqc.UploadMeasurementMetrics,0,mqc.OnReportResutl)

		topic:=mqc.DecodeResutlTopic
		log.Println("MQTTClient Subscribe topic:"+topic)
		client.Subscribe(topic,0,mqc.OnDecodeResutl)

		client.Subscribe(PINT_RESULT,0,mqc.OnResutl)
		client.Subscribe(ATTACH_RESULT,0,mqc.OnResutl)
		client.Subscribe(TCP_RESULT,0,mqc.OnResutl)
	}
}

func (mqc *MQTTClient) connectLostHandler(client mqtt.Client, err error){
	log.Println("MQTTClient connectLostHandler connect status: ",client.IsConnected(),err)
}

func (mqc *MQTTClient) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient messagePublishHandler topic: ",msg.Topic())
}

func (mqc *MQTTClient) reconnectingHandler(Client mqtt.Client,opts *mqtt.ClientOptions){
	log.Println("MQTTClient reconnectingHandler ")
}

func (mqc *MQTTClient) OnDecodeResutl(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient OnDecodeResutl ",msg.Topic(),string(msg.Payload()))

	if mqc.Handler != nil {
		mqc.Handler.HandleDecodeResult(string(msg.Payload()))
	}
}

func (mqc *MQTTClient) OnReportResutl(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient OnReportResutl ",msg.Topic(),string(msg.Payload()))

	if mqc.Handler != nil {
		mqc.ReportHandler.HandleReportResult(string(msg.Payload()))
	}
}

func (mqc *MQTTClient) OnResutl(Client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient OnResutl ",msg.Topic(),string(msg.Payload()))

	if mqc.Handler != nil {
		mqc.ReportHandler.HandleResult(string(msg.Payload()))
	}
}

func (mqc *MQTTClient)Publish(topic,content string)(int){
	if mqc.Client == nil {
		return common.ResultMqttClientError
	}
	log.Println("MQTTClient Publish topic:"+topic+" content:"+content)
	token:=mqc.Client.Publish(topic,0,false,content)
	token.Wait()
	return common.ResultSuccess
}

func (mqc *MQTTClient) Init(){
	mqc.getClient()
}