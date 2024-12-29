package testfile

import (
	"log"
)

type TestCmdSender struct {
	
}

func (dc *TestCmdSender)SendCmd(){
	log.Println("TestCmdSender.SendCmd")
}

func (dc *TestCmdSender)HasCmd()(bool){
	return true
}