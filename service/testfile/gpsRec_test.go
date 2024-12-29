package testfile

import (
	"testing"
	"rt_test_service/crv"
	"fmt"
)

var gpsStr="{\"Date\":\"241221\",\"LineCode\":\"\",\"LineName\":\"\",\"LineType\":\"\",\"SateliteCount\":11,\"UtcTime\":\"055904\",\"distance\":4.158,\"latitude\":39.85297,\"latitude_direction\":\"N\",\"longitude\":116.3502,\"longitude_direction\":\"E\",\"speed\":105.877}"
	

func _TestGetGPSRec(t *testing.T) {
	gpsRec:=GetGPSRec(gpsStr)

	if gpsRec==nil {
		t.Error("GetGPSRec failed")
		return
	}
}

func _TestSaveGPS(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:  "rt_test_tk_service",
		AppID:  "rb_test_tk",
	}

	SaveGPS(gpsStr,crvClient)
}

func TestGetGPSRecFromDB(t *testing.T) {
	crvClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:  "rt_test_tk_service",
		AppID:  "rb_test_tk",
	}

	gpsRec:=GetGPSRecFromDB("log1",crvClient)
	if gpsRec==nil {
		t.Error("GetGPSRec failed")
		return
	}

	fmt.Println("gpsRec:",gpsRec)
}