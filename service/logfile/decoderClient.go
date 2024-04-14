package logfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	URL_GET_STATUS     = "api/decoder/status_query"
	URL_DECODE_FILE    = "api/decoder/request"
)

type DecoderStatus struct {
	Status string `json:"status"`
	Details string `json:"details"`
}

type DecodeFileRequest struct {
	Logfiles *[]string `json:"logfiles"`
	Unixtime int64 `json:"unixtime"`
}

type DecodeFileResponse struct {
	Res string `json:"res"`
	Cause string `json:"cause"`
	ID int64 `json:"id"`
}

type DecoderClient struct {
	URL string
}

func (dc *DecoderClient) GetStatus()(*DecoderStatus, error) {
	log.Println("DecoderClient GetStatus start ...")
	req, err := http.NewRequest("GET", dc.URL+URL_GET_STATUS, nil)
	if err != nil {
		log.Println("DecoderClient GetStatus error", err)
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("DecoderClient GetStatus Do request error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("DecoderClient GetStatus StatusCode error", resp)
		return nil, errors.New("api "+URL_GET_STATUS+" returns a worng status code")
	}

	decoder := json.NewDecoder(resp.Body)
	res := &DecoderStatus{}
	err = decoder.Decode(res)
	if err != nil {
		log.Println("DecoderClient GetStatus result decode failed error", err.Error())
		return nil, err
	}

	resultJson, _ := json.Marshal(res)
	log.Println(string(resultJson))

	log.Println("end DecoderClient GetStatus success")
	return res, nil
}

func (dc *DecoderClient) DecodeFile(files *[]string)(*DecodeFileResponse, error) {
	logFile:=make([]string, len(*files))
	for i, file:=range *files{
		logFile[i]="http://192.168.4.111/qlogs/qlogs/09b59cee155449dca0de423adc4bc1a0/qlog_files/"+file
	}

	decodeFileRequest:=&DecodeFileRequest{
		Logfiles: &logFile,
		Unixtime: time.Now().Unix(),
	}
	reqJson, _ := json.Marshal(decodeFileRequest)

	log.Println("DecoderClient DecodeFile request", string(reqJson))

	postBody:=bytes.NewBuffer(reqJson)
	req, err := http.NewRequest("POST", dc.URL+URL_DECODE_FILE, postBody)
	if err != nil {
		log.Println("DecoderClient DecodeFile error", err)
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("DecoderClient DecodeFile Do request error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("DecoderClient DecodeFile StatusCode error", resp)
		return nil, errors.New("api "+URL_DECODE_FILE+" returns a worng status code")
	}

	decoder := json.NewDecoder(resp.Body)
	res := &DecodeFileResponse{}
	err = decoder.Decode(res)
	if err != nil {
		log.Println("DecoderClient DecodeFile result decode failed error", err.Error())
		return nil, err
	}

	resultJson, _ := json.Marshal(res)
	log.Println(string(resultJson))

	log.Println("end DecoderClient DecodeFile success")
	return res, nil	
}