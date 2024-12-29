package testfile

import (
	"encoding/json"
	"log"
	"rt_test_service/crv"
	"rt_test_service/common"
	"strconv"
)

type GPSRec struct {
	ID string `json:"id"`
	GlbDate string `json:"Date"`
	LineCode string `json:"LineCode"`
	LineName string `json:"LineName"`
	LineType string `json:"LineType"`
	SateliteCount int64 `json:"SateliteCount"`
	UtcTime string `json:"UtcTime"`
	Distance float64 `json:"distance"`
	Latitude float64 `json:"latitude"`
	LatitudeDirection string `json:"latitude_direction"`
	Longitude float64 `json:"longitude"`
	LongitudeDirection string `json:"longitude_direction"`
	Speed float64 `json:"speed"`
	Datetime string `json:"datetime"`
	GpsAntStatus string `json:"gps_ant_status"`
	GpsStatus string `json:"gps_status"`
	LogFile string `json:"log_file"`
}

func GetGPSRec(gpsStr string)(*GPSRec){
	gpsRec := GPSRec{}
	err := json.Unmarshal([]byte(gpsStr), &gpsRec)
	if err != nil {
		log.Println("GetGPSRec Unmarshal failed:", err)
		return nil
	}

	return &gpsRec
}

func SaveGPS(gpsStr string,crvClient *crv.CRVClient){
	log.Println("SaveGPS gpsStr:",gpsStr)
	gpsRec := GetGPSRec(gpsStr)
	if gpsRec == nil {
		return
	}

	SavePGSRec(gpsRec,crvClient)
}

func SavePGSRec(gpsRec *GPSRec,crvClient *crv.CRVClient){
	grsRecMap := map[string]interface{}{
		"glb_date":gpsRec.GlbDate,
		"line_code":gpsRec.LineCode,
		"line_name":gpsRec.LineName,
		"line_type":gpsRec.LineType,
		"satelite_count":strconv.FormatInt(gpsRec.SateliteCount, 10),
		"utctime":gpsRec.UtcTime,
		"distance":gpsRec.Distance,
		"latitude":gpsRec.Latitude,
		"latitude_direction":gpsRec.LatitudeDirection,
		"longitude":gpsRec.Longitude,
		"longitude_direction":gpsRec.LongitudeDirection,
		"speed":gpsRec.Speed,
		"datetime":gpsRec.Datetime,
		"gps_ant_status":gpsRec.GpsAntStatus,
		"gps_status":gpsRec.GpsStatus,
		"log_file":gpsRec.LogFile,
		"_save_type": "create",
	}

	commonRep := crv.CommonReq{
		ModelID: "rt_gps_rec",
		List:    &[]map[string]interface{}{grsRecMap},
	}

	crvClient.Save(&commonRep, "")
}

func GetGPSRecFromDB(logFile string,crvClient *crv.CRVClient)(*[]GPSRec){
	commonRep := crv.CommonReq{
		ModelID: "rt_gps_rec",
		Filter:&map[string]interface{}{
			"log_file":map[string]interface{}{
				"Op.eq":logFile,
			},
		},
		Fields:  &[]map[string]interface{}{
			{"field": "id"},		
			{"field": "glb_date"},
			{"field": "line_code"},
			{"field": "line_name"},
			{"field": "line_type"},
			{"field": "satelite_count"},
			{"field": "utctime"},
			{"field": "distance"},
			{"field": "latitude"},
			{"field": "latitude_direction"},
			{"field": "longitude"},
			{"field": "longitude_direction"},
			{"field": "speed"},
			{"field": "datetime"},
			{"field": "gps_ant_status"},
			{"field": "gps_status"},
			{"field": "log_file"},
		},
		Pagination:&crv.Pagination{
			Current:1,
			PageSize:100000,
		},
	}

	rsp, commonErr := crvClient.Query(&commonRep, "")
	if commonErr != common.ResultSuccess {
		return nil
	}

	if rsp.Error == true {
		log.Println("GetGPSRec error:", rsp.ErrorCode, rsp.Message)
		return nil
	}

	resLst, ok := rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetGPSRec error: no list in rsp.")
		return nil
	}

	if len(resLst) == 0 {
		log.Println("GetGPSRec error: no list in rsp.")
		return nil
	}

	gpsRecList:=[]GPSRec{}
	
	for _, res := range resLst {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			log.Println("GetGPSRec error: res is not map.")
			continue
		}

		sateliteStr:=resMap["satelite_count"].(string)
		var sateliteCount int64
		i, err := strconv.ParseInt(sateliteStr, 10, 64)
		if err == nil {
			sateliteCount=i
		} else {
			log.Println("satelite_count is not int. sateliteStr:"+sateliteStr)
		}

		distanceStr:=resMap["distance"].(string)
		distance:=0.0
		f, err := strconv.ParseFloat(distanceStr, 64)
		if err == nil {
			distance=f
		} else {
			log.Println("distance is not float. distanceStr:"+distanceStr)
		}

		latitudeStr:=resMap["latitude"].(string)
		latitude:=0.0
		f, err = strconv.ParseFloat(latitudeStr, 64)
		if err == nil {
			latitude=f
		} else {
			log.Println("latitude is not float. latitudeStr:"+latitudeStr)
		}

		longitudeStr:=resMap["longitude"].(string)
		longitude:=0.0
		f, err = strconv.ParseFloat(longitudeStr, 64)
		if err == nil {
			longitude=f
		} else {
			log.Println("longitude is not float. longitudeStr:"+longitudeStr)
		}

		speedStr:=resMap["speed"].(string)
		speed:=0.0
		f, err = strconv.ParseFloat(speedStr, 64)
		if err == nil {
			speed=f
		} else {
			log.Println("speed is not float. speedStr:"+speedStr)
		}

		gpsRec := GPSRec{
			ID: resMap["id"].(string),
			GlbDate: resMap["glb_date"].(string),
			LineCode: resMap["line_code"].(string),
			LineName: resMap["line_name"].(string),
			LineType: resMap["line_type"].(string),
			SateliteCount:  sateliteCount,
			UtcTime: resMap["utctime"].(string),
			Distance: distance,
			Latitude: latitude,
			LatitudeDirection: resMap["latitude_direction"].(string),
			Longitude: longitude,
			LongitudeDirection: resMap["longitude_direction"].(string),
			Speed: speed,
			Datetime: resMap["datetime"].(string),
			GpsAntStatus: resMap["gps_ant_status"].(string),
			GpsStatus: resMap["gps_status"].(string),
			LogFile: resMap["log_file"].(string),
		}

		gpsRecList = append(gpsRecList, gpsRec)
	}

	return &gpsRecList
}