package device

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	URL_GET_IMSI     = "api/ue/get_imsi"
	URL_DIAL_QUERY   = "api/ue/dial_query"
	URL_DIAL_TRIGGER = "api/ue/dial_trigger"
	URL_DEVICE_REBOOT = "api/ue/reset"
	URL_ATTACH       = "api/ue/attach"
	URL_ATTACH_QUERY = "api/ue/attach_query"
	URL_DETACH       = "api/ue/detach"
	URL_GET_RAT = "api/ue/get_rat_status"
	URL_SET_RAT = "api/ue/nr5g"
)

type DeviceClient struct {
	ServerUrl string
}

func (dc *DeviceClient) Get(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", dc.ServerUrl+url, nil)
	if err != nil {
		log.Println("DeviceClient Get error", err)
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println("DeviceClient Get Do request error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("DeviceClient Get StatusCode error", resp)
		return nil, errors.New("api returns a worng status code")
	}

	decoder := json.NewDecoder(resp.Body)
	res := map[string]interface{}{}
	err = decoder.Decode(&res)
	if err != nil {
		log.Println("DeviceClient Get result decode failed error", err.Error())
		return nil, err
	}

	resultJson, _ := json.Marshal(&res)
	log.Println(string(resultJson))

	log.Println("end DeviceClient Get success")
	return res, nil
}

func (dc *DeviceClient) GetImsi() (map[string]interface{}, error) {
	return dc.Get(URL_GET_IMSI)
}

func (dc *DeviceClient) GetDialQuery() (map[string]interface{}, error) {
	return dc.Get(URL_DIAL_QUERY)
}

func (dc *DeviceClient) DialTrigger() (map[string]interface{}, error) {
	return dc.Get(URL_DIAL_TRIGGER)
}

func (dc *DeviceClient) DeviceReboot() (map[string]interface{}, error) {
	return dc.Get(URL_DEVICE_REBOOT)
}

func (dc *DeviceClient) Attach() (map[string]interface{}, error) {
	return dc.Get(URL_ATTACH)
}

func (dc *DeviceClient) AttachQuery() (map[string]interface{}, error) {
	return dc.Get(URL_ATTACH_QUERY)
}

func (dc *DeviceClient) Detach() (map[string]interface{}, error) {
	return dc.Get(URL_DETACH)
}

func (dc *DeviceClient)GetRAT()(interface{}, error){
	return dc.Get(URL_GET_RAT)
}

func (dc *DeviceClient)SetRAT(rat string)(interface{}, error){
	return dc.Get(URL_SET_RAT)
}
