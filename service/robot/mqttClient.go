package robot

import (
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"rt_test_service/common"
	"github.com/rs/xid"
)

const (
	MSG_TYPE_DIAG="Diag"
	MSG_TYPE_EVENT="Event"
	MSG_TYPE_SIGNAL="SignalFilter"
)

type eventHandler interface {
	DealDeviceHeartbeat(deviceID,vin string)
	DealDiagResponse(deviceID string)
	DealEventResponse(deviceID string)
	DealSignalResponse(deviceID string)
}

type RobotMQTTClient struct {
	Broker string
	User string
	Password string
	Client mqtt.Client
}

func (mqc *RobotMQTTClient)getClientID()(string){
	return xid.New().String()
}

func (mqc *RobotMQTTClient) getClient()(mqtt.Client){
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqc.Broker)
	opts.SetClientID(mqc.getClientID())
	opts.SetUsername(mqc.User)
	opts.SetPassword(mqc.Password)
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(mqc.messagePublishHandler)
	opts.OnConnect = mqc.connectHandler
	opts.OnConnectionLost = mqc.connectLostHandler
	opts.OnReconnecting = mqc.reconnectingHandler

	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error)
		return nil
	}
	return client
}

func (mqc *RobotMQTTClient) connectHandler(client mqtt.Client){
	log.Println("MQTTClient connectHandler connect status: ",client.IsConnected())
}

func (mqc *RobotMQTTClient) connectLostHandler(client mqtt.Client, err error){
	log.Println("MQTTClient connectLostHandler connect status: ",client.IsConnected(),err)
}

func (mqc *RobotMQTTClient) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("MQTTClient messagePublishHandler topic: ",msg.Topic())
}

func (mqc *RobotMQTTClient) reconnectingHandler(Client mqtt.Client,opts *mqtt.ClientOptions){
	log.Println("MQTTClient reconnectingHandler ")
}

func (mqc *RobotMQTTClient)onRobotStatus(Client mqtt.Client, msg mqtt.Message){
	log.Println("RobotMQTTClient onRobotStatus ",msg.Topic())
	strTopic:=msg.Topic()
	log.Println("RobotMQTTClient onRobotStatus strTopic ",strTopic)
	log.Println("RobotMQTTClient onRobotStatus msg ",string(msg.Payload()))
}

func (mqc *RobotMQTTClient)Publish(topic,content string)(int){
	if mqc.Client == nil {
		return common.ResultMqttClientError
	}
	log.Println("MQTTClient Publish topic:"+topic+" content:"+content)
	token:=mqc.Client.Publish(topic,0,false,content)
	token.Wait()
	return common.ResultSuccess
}

func (mqc *RobotMQTTClient) Init(){
	mqc.Client=mqc.getClient()
	mqc.Client.Subscribe("XJPHrobotStatus",0,mqc.onRobotStatus)
	//mqc.Client.Subscribe(mqc.DiagResponseTopic,0,mqc.onDiagResponse)
}