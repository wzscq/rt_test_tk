package robot

import (
	"rt_test_service/common"
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"bytes"
)

type RobotStatus struct {
	RobotId string `json:"robotId"`
	RobotName string `json:"robotName"`
	Electricity string `json:"electricity"`
	MapId string `json:"mapId"`
	MapName string `json:"mapName"`
	TaskId string `json:"taskId"`
	TaskName string `json:"taskName"`
	CurrentTask string `json:"currentTask"`
	RobotStatus string `json:"robotStatus"`
	PixelX string `json:"pixelX"`
	PixelY string `json:"pixelY"`
	PixelTheta string `json:"pixelTheta"`
	Exception string `json:"exception"`
	DatetimeSend string `json:"dateTimeSend"`
	StateCode string `json:"stateCode"`
}

type RobotInfo struct {
	RobotId string `json:"robotId"`
	RobotName string `json:"robotName"`
	LongitudeLatitude string `json:"longitudeLatitude"`
	Position string `json:"position"`
	Electricity string `json:"electricity"`
	NewType string `json:"newType"`
	Online string `json:"online"`
	IsOnlineStatus bool `json:"isOnlineStatus"`
}

type SlotInfo struct {
	Country string `json:"country"`
	SlotID string `json:"slot_id"`
	IMEI string `json:"imei"`
	PLMN string `json:"plmn"`
	IMSI string `json:"imsi"`
	Short string `json:"short"`
	Operator string `json:"operator"`
}

type UEInfo struct {
	RilImpl string `json:"ril_impl"`
  ChipManufacturer string `json:"chip_manufacturer"`
  MarketName string `json:"marketname"`
  ProductManufacturer string `json:"product_manufacturer"`
	Slot []SlotInfo `json:"slot"`
	SN string `json:"sn"`
  Baseband string `json:"baseband"`
}

type EquipmentStatus struct {
	Phase string `json:"phase"`
	RobotId string `json:"robot_id"`
	UEDetails []UEInfo `json:"ue_details"`
	HostStatus string `json:"host_status"`
	HostId string `json:"host_id"`
}

type getTestEquipmentStatusRsp struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"code"`
	Result *EquipmentStatus `json:"result"`
	Timestamp int64 `json:"timestamp"`
}

type sendTaskRsp struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"code"`
	Result interface{} `json:"result"`
	Timestamp int64 `json:"timestamp"`
}

type getCurrentRobotStatusRsp struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"code"`
	Result *RobotStatus `json:"result"`
	Timestamp int64 `json:"timestamp"`
}

type oauthReq struct {
	AppKey string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

type oauthResult struct {
	ExpiresIn string `json:"expiresIn"`
	Token string `json:"token"`
}

type getRobotListRsp struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"code"`
	Result []RobotInfo `json:"result"`
	Timestamp int64 `json:"timestamp"`
}

type oauthRsp struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"code"`
	Result *oauthResult `json:"result"`
	Timestamp int64 `json:"timestamp"`
}

type RobotClient struct {
	Conf *common.RobotClientConf
}

func (rc *RobotClient) Oauth()(*oauthRsp, error) {
	log.Println("RobotClient Oauth ...")
	rep:=oauthReq{
		AppKey:rc.Conf.Oauth.AppKey,
		AppSecret:rc.Conf.Oauth.AppSecret,
	}

	params := url.Values{}
	params.Set("appKey", rc.Conf.Oauth.AppKey)
	params.Add("appSecret", rc.Conf.Oauth.AppSecret)
	
	url:=rc.Conf.Oauth.URL+"?"+params.Encode()

	postJson,_:=json.Marshal(rep)
	postBody:=bytes.NewBuffer(postJson)
	resp,err:=http.Post(url,"application/json",postBody)

	if err != nil {
		log.Println("RobotClient Oauth error",err)
		return nil,err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println("RobotClient Oauth error",resp)
		return nil,err
	}

	decoder := json.NewDecoder(resp.Body)
	rsp:=oauthRsp{}
	err = decoder.Decode(&rsp)
	if err != nil {
		log.Println("RobotClient Oauth rsp decode failed [Err:%s]", err.Error())
		return nil,err
	}

	log.Println("RobotClient Oauth success with token:",rsp.Result.Token)
	return &rsp,nil
}

func (rc *RobotClient) GetRobotList(token string)(*getRobotListRsp, error) {
	req,err:=http.NewRequest("GET",rc.Conf.GetRobotList.URL,nil)
	if err != nil {
		log.Println("RobotClient GetRobotList NewRequest error",err)
		return nil,err
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Access-Token",token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("RobotClient GetRobotList Do request error",err)
		return nil,err
	}
	defer resp.Body.Close()


	decoder := json.NewDecoder(resp.Body)
	rsp:=getRobotListRsp{}
	err = decoder.Decode(&rsp)
	if err != nil {
		log.Println("RobotClient GetRobotList rsp decode failed [Err:%s]", err.Error())
		return nil,err
	}

	log.Println("RobotClient GetRobotList success with result:",rsp)
	return &rsp,nil
}

func (rc *RobotClient) GetCurrentRobotStatus(token,robotID string)(*getCurrentRobotStatusRsp, error) {
	
	params := url.Values{}
	params.Set("robotId", robotID)
	
	url:=rc.Conf.GetCurrentRobotStatus.URL+"?"+params.Encode()

	req,err:=http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("RobotClient GetCurrentRobotStatus NewRequest error",err)
		return nil,err
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Access-Token",token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("RobotClient GetCurrentRobotStatus Do request error",err)
		return nil,err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	rsp:=getCurrentRobotStatusRsp{}
	err = decoder.Decode(&rsp)
	if err != nil {
		log.Println("RobotClient GetCurrentRobotStatus rsp decode failed [Err:%s]", err.Error())
		return nil,err
	}

	log.Println("RobotClient GetCurrentRobotStatus success with result:",rsp)
	return &rsp,nil
}

func (rc *RobotClient) GetTestEquipmentStatus(token,robotID string)(*getTestEquipmentStatusRsp, error) {
	params := url.Values{}
	params.Set("robotId", robotID)
	
	url:=rc.Conf.GetTestEquipmentStatus.URL+"?"+params.Encode()

	req,err:=http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("RobotClient GetTestEquipmentStatus NewRequest error",err)
		return nil,err
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Access-Token",token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("RobotClient GetTestEquipmentStatus Do request error",err)
		return nil,err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	rsp:=getTestEquipmentStatusRsp{}
	err = decoder.Decode(&rsp)
	if err != nil {
		log.Println("RobotClient GetTestEquipmentStatus rsp decode failed [Err:%s]", err.Error())
		return nil,err
	}

	log.Println("RobotClient GetTestEquipmentStatus success with result:",rsp)
	return &rsp,nil
}

func (rc *RobotClient)SendTask(token string,task map[string]interface{})(*sendTaskRsp,error){
	url:=rc.Conf.SendTask.URL

	postJson,_:=json.Marshal(task)

	log.Println("RobotClient SendTask with json:",string(postJson))

	postBody:=bytes.NewBuffer(postJson)

	req,err:=http.NewRequest("POST",url,postBody)
	if err != nil {
		log.Println("RobotClient SendTask NewRequest error",err)
		return nil,err
	}

	req.Header.Set("Content-Type","application/json")
	req.Header.Set("X-Access-Token",token)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("RobotClient SendTask Do request error",err)
		return nil,err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	rsp:=sendTaskRsp{}
	err = decoder.Decode(&rsp)
	if err != nil {
		log.Println("RobotClient SendTask rsp decode failed [Err:%s]", err.Error())
		return nil,err
	}

	log.Println("RobotClient SendTask success with result:",rsp)
	return &rsp,nil
}

