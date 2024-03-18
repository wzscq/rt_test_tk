package logfile

import (
	"rt_test_service/crv"
	"encoding/json"
	"log"
)

type DecodeResultHandler struct {
	CRVClient *crv.CRVClient
}

type DecodeStatus struct {
	Status string `json:"status"`
	CurrentPhase string `json:"current_phase"`
	DecodeFiles string `json:"decode_files"`
	DecodeID int64 `json:"id"`
}

func (drh *DecodeResultHandler) HandleDecodeResult(result string) {
	//decode result
	decodeStatus:= DecodeStatus{}
	err := json.Unmarshal([]byte(result), &decodeStatus)
	//update db
	if err != nil {
		log.Println("HandleDecodeResult error: ", err)
		return
	}

	UpdateDecodeStatus(&decodeStatus, drh.CRVClient,"")
}