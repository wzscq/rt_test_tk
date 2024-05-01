package logfile

import (
	"testing"
	"fmt"
	"rt_test_service/crv"
)

func _TestGetStatus(t *testing.T){
	decoderClient:=DecoderClient{
		URL:"http://182.42.81.6:5000/",
	}

	status,err:=decoderClient.GetStatus()
	if err!=nil{
		t.Error("GetStatus error")
		return
	}

	fmt.Println(status)
}

func _TestDecodeFile(t *testing.T){
	decoderClient:=DecoderClient{
		URL:"http://182.42.81.6:5000/",
	}

	rsp,err:=decoderClient.DecodeFile(&[]string{"test.log"})
	if err!=nil{
		t.Error("DecodeFile error")
		return
	}

	fmt.Println(rsp)
}

func _TestDecodeLogFile(t *testing.T){
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	decoderClient:=DecoderClient{
		URL:"http://182.42.81.6:5000/",
	}

	err:=DecodeLogFile(&[]string{"test.log"},&decoderClient,crvClient,"rt_test_tk_service")
	if err!=nil{
		t.Error("DecodeLogFile error")
		return
	}
}

func TestGetDecodingTaskCount(t *testing.T){
	crvClient := &crv.CRVClient{
		Server: "http://localhost:8200",
		Token:  "rt_test_tk_service",
		AppID:  "",
	}

	count,err:=GetDecodingTaskCount(crvClient,"rt_test_tk_service")
	if err!=nil{
		t.Error("GetDecodingTaskCount error")
		return
	}

	fmt.Println(count)
}