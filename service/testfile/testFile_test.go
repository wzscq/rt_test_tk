package testfile

import (
    "testing"
	"time"
	"strconv"
	"log"
	"encoding/json"
	"rt_test_service/crv"
	"fmt"
)

var reportData=map[string]interface{}{
	"data": []map[string]interface{}{
			{
				"radio": map[string]interface{}{
					"measures_common": map[string]interface{}{
						"Current Network Type": "",
						"IMSI": "460011895631209",
						"IMEI": "864451045593183",
						"imsi": "460011895631209",
						"imei": "864451045593183",
						"duplex_mode": "",
						"net_Tye": "NR_SA",
						"plmn": "46001",
						"operator": "\u4e2d\u56fd\u8054\u901a",
						"operator_country": "cn",
						"operator_short": "CHN-UNICOM",
						"phone_name": "Redmi K30 5G Speed",
					},
					"measures_nr": map[string]interface{}{
						"NR PHY Throughput UL(CA Total)": "231528",
						"NR PHY Throughput DL(CA Total)": "3723656",
						"NR C-RNTI": "",
						"NR PUSCH Initial BLER": "7.69",
						"NR PUSCH BLER": "7.14",
						"NR PDSCH Initial BLER": "2.25",
						"NR PDSCH BLER": "2.2",
						"NR UL Avg MCS": "24.21",
						"NR DL Avg MCS": "10.76",
						"NR High Modulation UL/s": "3",
						"NR High Modulation DL/s": "4",
						"NR PUSCH TxPower": "15",
						"NR PUCCH TxPower": "19",
						"NR PRACH TxPower": "",
						"NR PDCCH UL GrantCount": "45",
						"NR PDCCH DL GrantCount": "99",
						"SS-SINR": "11.37",
						"SS-RSRP": "-84.59",
						"NR Network State": "NR Connected",
						"UE Category": "",
						"NR Slot Config UL Total": "",
						"NR Slot Config DL Total": "",
						"NR CSI RS Period": "",
						"NR SSB Period": "",
						"SSB Config Beam Num": "",
						"NR PCI": "53",
						"SSB ARFCN": "627264",
						"NR SSB Frequency(MHz)": "3408960",
						"NR WB CQI": "11.31",
						"NCI": "",
						"NR Bandwidth(MHz)": "",
						"NR Band": "",
						"NR TAC": "",
					},
				},
				"event": []map[string]interface{}{
					{
						"Logfile": "testRecoder_0818-212920",
						"pointIndex": 204,
						"EventIndex": 11,
						"EventTime": "2023-08-18 21:29:29.121",
						"EventCode": "0x8D",
						"name": "FTP Download TCPSlow",
					},
				},
				"msg": []map[string]interface{}{
					{
						"Logfile": "testRecoder_0818-212920",
						"pointIndex": 430,
						"MsgTime": "2023-08-18 21:29:33.894",
						"MsgCode": "0x50000200",
						"name": "NR->MeasurementReport",
					},
				},
				"time": 1692365376,
				"case_progress": map[string]interface{}{
					"session_id": 161920005569,
					"full_path": "d:\\test_recoder\\testReport_001_161920005569\\progress_460011895631209.csv",
					"imsi": "460011895631209",
					"DevId": "28e79a383c89153af87780ec745e2481",
					"Status": "FTPDownload",
					"Times": "1/1",
					"Progress": "8/30s",
					"FailTimes": 0,
					"OtherInfo": "Inst:2,729.30 kbps  Avg:2,799.36 kbps SuccRatio:100%(1/1)",
				},
			},
	},
	"robot_info": map[string]interface{}{
			"id": 1,
			"robot_id": "2bee174b7d7c36e9b98bd8772e66af5e",
			"map_id": "9ec5a9c61a374be4b0570077d7d6dd05",
			"pixel_x": 616.645741420779,
			"pixel_y": 285.72210123115497,
			"pixel_theta": -110.04722431915482,
			"record": "2023-08-18 21:27:05",
	},
	"pcTime": 1692365376,
}

func _TestCreateFile(t *testing.T) {
	outPath:="../localcache/"
	deviceID:="device1"
	timeStamp:=fmt.Sprintf("%d",time.Now().Unix())
	tf:=GetTestFile(outPath,deviceID,timeStamp)
	if tf==nil {
		t.Error("GetTestFile failed")
		return
	}

	for i:=0;i<10;i++ {
		lineContent:="line"+strconv.Itoa(i)
		tf.WriteLine(lineContent)
	}
	
	tf.Close()
}


func _TestTestFilePool(t *testing.T){
	
	crvClient:=&crv.CRVClient{
		Server:"http://127.0.0.1:8200",
		Token:"rt_test_service",
		AppID:"rt_test",
	}
	
	//reportData转换为JSON字符串
	reportDataJson,_:=json.Marshal(reportData)
	//创建TestFilePool
	tfp:=InitTestFilePool("../localcache/","3s",crvClient)
	//写入数据
	tfp.WriteDeviceTestLine("device2",string(reportDataJson))
	//等待5秒
	for i:=0;i<5;i++ {
		time.Sleep(2*time.Second)
		tfp.WriteDeviceTestLine("device2",string(reportDataJson))
	}
	time.Sleep(8*time.Second)
}

func _TestGetFilePoints(t *testing.T){
	outPath:="../localcache"
	deviceID:="BFEBFBFF000806EC"
	timestamp:="1694409161"

	tf:=GetTestFile(outPath,deviceID,timestamp)
	if tf==nil {
		t.Error("TestFileController GetPoints file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	indicator:=Indicator{
		ExtractPath:"radio.measures_lte.RSRP",
		ID:"1",
		Name:"test",
		Legend:IndicatorLegend{
			ModelID:"model",
			Total:3,
			List:[]IndicatorLegendItem{
				IndicatorLegendItem{
					ID:"id",
					SN:"1",
					Start:"",
					End:"-80.75",
					RGB:"#110000",
				},
				IndicatorLegendItem{
					ID:"id",
					SN:"2",
					Start:"-80.75",
					End:"-80.50",
					RGB:"#220000",
				},
				IndicatorLegendItem{
					ID:"id",
					SN:"2",
					Start:"-80.50",
					End:"-80.24",
					RGB:"#330000",
				},
				IndicatorLegendItem{
					ID:"id",
					SN:"2",
					Start:"-80.25",
					End:"",
					RGB:"#440000",
				},
			},
		},
	}

	//获取文件内容
	points:=tf.GetPoints(indicator)

	pointsJson,_:=json.Marshal(points)	

	log.Println(string(pointsJson))
}

func TestGetLineTimeStamp(t *testing.T){
	line:=`{"data": [{"radio": {"measures_common": {"imsi": "460001933146754", "imei": "864451045593183", "duplex_mode": "", "net_Tye": "LTE", "plmn": "46000", "operator": "\u4e2d\u56fd\u79fb\u52a8", "operator_country": "cn", "operator_short": "CMCC", "phone_name": "Redmi K30 5G Speed"}, "measures_lte": {"UE Category": "", "TAC": "", "Cell ID": "", "Band": "3", "EARFCN DL": "1309", "PCI": "12", "Bandwidth DL(MHz)": "", "TM": "", "Frequency DL(MHz)": "1815.9", "SubFrame Assign Type": "", "Special SubFrame Patterns": "", "LTE Network State": "", "RRC Protocol": "17", "M-TMSI": "", "C-RNTI": "", "CodeWord Number": "1", "RSRQ": "-6.81", "RSRP": "-80.56", "SINR": "18.2", "RSSI": "-53.75", "Pathloss": "", "PRACH TxPower": "", "PUCCH TxPower": "", "PUSCH TxPower": "", "PDSCH RB Count/s": "-1200", "PDSCH Scheduled SubFN Count /s": "1", "PDSCH Scheduled RB Count /slot": "6", "PUSCH RB Count/s": "", "PUSCH Scheduled SubFN Count /s": "", "PUSCH Scheduled RB Count /slot": "", "MCS Average DL /s": "0", "MCS Average UL /s": "", "Rank1 CQI": "", "Rank2 CQI Code0": "", "Rank2 CQI Code1": "", "PDSCH BLER": "0", "PUSCH BLER": "", "PHY Throughput DL(CA Total)": "", "PHY Throughput UL(CA Total)": ""}, "measures_nr": {"NR PHY Throughput UL(CA Total)": "", "NR PHY Throughput DL(CA Total)": "", "NR C-RNTI": "", "NR PUSCH Initial BLER": "", "NR PUSCH BLER": "", "NR PDSCH Initial BLER": "", "NR PDSCH BLER": "", "NR UL Avg MCS": "", "NR DL Avg MCS": "", "NR High Modulation UL/s": "", "NR High Modulation DL/s": "", "NR PUSCH TxPower": "", "NR PUCCH TxPower": "", "NR PRACH TxPower": "", "NR PDCCH UL GrantCount": "", "NR PDCCH DL GrantCount": "", "SS-SINR": "", "SS-RSRP": "", "NR Network State": "", "UE Category": "", "NR Slot Config UL Total": "", "NR Slot Config DL Total": "", "NR CSI RS Period": "", "NR SSB Period": "", "SSB Config Beam Num": "", "NR PCI": "", "SSB ARFCN": "", "NR SSB Frequency(MHz)": "", "NR WB CQI": "", "NCI": "", "NR Bandwidth(MHz)": "", "NR Band": "", "NR TAC": ""}}, "event": [{"Logfile": "testRecoder_0911-131237", "pointIndex": 7, "EventIndex": 1, "EventTime": "2023-09-11 13:12:40.245", "EventCode": "0x19F", "name": "Test Describe"}, {"Logfile": "testRecoder_0911-131237", "pointIndex": 7, "EventIndex": 2, "EventTime": "2023-09-11 13:12:40.248", "EventCode": "0x19C", "name": "Test plan execute Request Event"}, {"Logfile": "testRecoder_0911-131237", "pointIndex": 14, "EventIndex": 3, "EventTime": "2023-09-11 13:12:40.867", "EventCode": "0xA", "name": "Voice Dial"}], "msg": [], "time": 1694409161, "case_progress": {"session_id": 161921005832, "full_path": "d:\\test_recoder\\testReport_001_161921005832\\progress_460001933146754.csv", "imsi": "460001933146754", "DevId": "3046021acf2686d9ab90d5ae8e696ed2", "Status": "Waiting", "Times": "", "Progress": "", "FailTimes": 0, "OtherInfo": ""}}], "robot_info": {"id": 1, "robot_id": "2bee174b7d7c36e9b98bd8772e66af5e", "map_id": "9ec5a9c61a374be4b0570077d7d6dd05", "pixel_x": 616.645741420779, "pixel_y": 285.72210123115497, "pixel_theta": -110.04722431915482, "record": "2023-09-11 12:02:12"}, "pcTime": 1694409161}`
	timeStamp:=GetLineTimeStamp(line)
	if timeStamp=="" {
		t.Error("GetLineTimeStamp failed")
		return
	}

	if timeStamp!="161921005832" {
		t.Error("GetLineTimeStamp failed")
		return
	}
	log.Println(timeStamp)
}	