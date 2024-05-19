package device

import (
	"rt_test_service/mqtt"
)

type CmdSender interface {
	HasCmd()(bool)
	SendCmd()
}

type CmdSenderRunner interface {
	SetCmdSender(CmdSender)
}

type DefaultCmdSender struct {
	CMD *TestCommand
	MQTTClient *mqtt.MQTTClient
	Topic string
}

func (dc *DefaultCmdSender)SendCmd(){
	SendTestCase(dc.CMD,dc.MQTTClient,dc.Topic)
}

func (dc *DefaultCmdSender)HasCmd()(bool){
	return dc.CMD!=nil
}