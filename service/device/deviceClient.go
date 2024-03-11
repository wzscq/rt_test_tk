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
