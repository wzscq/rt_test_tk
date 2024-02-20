package robot

import (
	"rt_test_service/common"
	"rt_test_service/crv"
	"log"
)

var MODELID_DEVICE="rt_device"

var	queryDeviceFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
	{
		"field": "device_ues",
		"fieldType":"one2many",
		"relatedModelID":"rt_device_ue",
		"relatedField":"device_id",
		"fields":[]map[string]interface{}{
			{"field": "id"},
			{"field": "sn"},
			{"field": "version"},
			{
				"field": "ue_imeis",
				"fieldType":"one2many",
				"relatedModelID":"rt_ue_imei",
				"relatedField":"device_ue_id",
				"fields":[]map[string]interface{}{
					{"field": "id"},
					{"field": "imei"},
					{"field": "imsi"},
					{"field": "version"},
				},
		 	},
		},
  },
}

func GetDevice(crvClient *crv.CRVClient,deviceID string,token string)(map[string]interface{},*common.CommonRsp){

	commonRep:=crv.CommonReq{
		ModelID:MODELID_DEVICE,
		Fields:&queryDeviceFields,
		Filter:&map[string]interface{}{"id":deviceID},
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return nil,common.CreateResponse(common.CreateError(commonErr,nil),nil)
	}

	if rsp.Error == true {
		log.Println("GetDevice error:",rsp.ErrorCode,rsp.Message)
		return nil,rsp
	}

	deviceList,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetDevice error: no list in rsp.")
		return nil,nil
	}

	if len(deviceList) == 0 {
		log.Println("GetDevice error: no device found.")
		return nil,nil
	}

	deviceItem,ok:=deviceList[0].(map[string]interface{})
	if !ok {
		log.Println("GetDevice error: can not convert item to map.")
		return nil,nil
	}

	return deviceItem,nil
}

func UpdateRobotEquipment(crvClient *crv.CRVClient,equipmentStatus *EquipmentStatus,token string)(*common.CommonRsp){
	//获取已经存在的机器人设备列表
	deviceItem,commonRsp:=GetDevice(crvClient,equipmentStatus.RobotId,token)
	if commonRsp != nil {
		return commonRsp
	}

	//如果设备不存在则创建
	if deviceItem == nil {
		return CreateDevice(crvClient,equipmentStatus,token)
	}
	//如果设备已经存在则更新
	return UpdateDevice(crvClient,equipmentStatus,deviceItem,token)
}

func CreateDevice(crvClient *crv.CRVClient,equipmentStatus *EquipmentStatus,token string)(*common.CommonRsp){
	device_ues:=[]map[string]interface{}{}
	for _,ueItem:= range equipmentStatus.UEDetails {
		deviceUe:=CreateDeviceUe(&ueItem,equipmentStatus.RobotId)
		device_ues=append(device_ues,deviceUe)
	}

	device:=map[string]interface{}{
		"id":equipmentStatus.RobotId,
		"host_status":equipmentStatus.HostStatus,
		"host_id":equipmentStatus.HostId,
		"phase":equipmentStatus.Phase,
		"_save_type":"create",
		"device_ues":map[string]interface{}{
			"fieldType":"one2many",
			"modelID":"rt_device_ue",
			"relatedField":"device_id",
			"list":device_ues,
		},
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_DEVICE,
		List:&[]map[string]interface{}{device},
	}

	rsp,commonErr:=crvClient.Save(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		return rsp
	}

	if rsp.Error == true {
		log.Println("CreateDevice error:",rsp.ErrorCode,rsp.Message)
		return rsp
	}

	return nil
}

func UpdateDevice(crvClient *crv.CRVClient,equipmentStatus *EquipmentStatus,deviceItem map[string]interface{},token string)(*common.CommonRsp){
	
	device_ues:=GetUpdateDeviceUes(equipmentStatus,deviceItem)

	device:=map[string]interface{}{
		"id":deviceItem["id"],
		"version":deviceItem["version"],
		"_save_type":"update",
		"host_status":equipmentStatus.HostStatus,
		"host_id":equipmentStatus.HostId,
		"phase":equipmentStatus.Phase,
		"device_ues":map[string]interface{}{
			"fieldType":"one2many",
			"modelID":"rt_device_ue",
			"relatedField":"device_id",
			"list":device_ues,
		},
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_DEVICE,
		List:&[]map[string]interface{}{device},
	}

	rsp,commonErr:=crvClient.Save(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		return rsp
	}

	if rsp.Error == true {
		log.Println("UpdateDevice error:",rsp.ErrorCode,rsp.Message)
		return rsp
	}

	return nil
}

func GetSameUe(sn string,deviceItem map[string]interface{})(map[string]interface{}){
	devUeMap,ok:=deviceItem["device_ues"].(map[string]interface{})
	if !ok {
		return nil
	}

	devUeList,ok:=devUeMap["list"]
	if !ok || devUeList == nil {
		return nil
	}

	for _,ueItem:= range devUeList.([]interface{}) {
		if ueItem.(map[string]interface{})["sn"] == sn {
			return ueItem.(map[string]interface{})
		}
	}
	return nil
}

func CreateDeviceUeImei(slot *SlotInfo,deviceID string)(map[string]interface{}){
	return map[string]interface{}{
		"device_id":deviceID,
		"sn":slot.SlotID,
		"imei":slot.IMEI,
		"imsi":slot.IMSI,
		"operator":slot.Operator,
		"country":slot.Country,
		"short":slot.Short,
		"plmn":slot.PLMN,
		"_save_type":"create",
	}
}

func CreateDeviceUe(ueItem *UEInfo,deviceID string)(map[string]interface{}){
		ueImeis:=[]map[string]interface{}{}
		for _,slot:= range ueItem.Slot {
			ueImei:=CreateDeviceUeImei(&slot,deviceID)
			ueImeis=append(ueImeis,ueImei)
		}

		deviceUe:=map[string]interface{}{
			"device_id":deviceID,
			"chip_manufacturer":ueItem.ChipManufacturer,
			"product_manufacturer":ueItem.ProductManufacturer,
			"marketname":ueItem.MarketName,
			//"phone_type":ueItem.PhoneType,
			//"ue_status":ueItem.UeStatus,
			//"network_type":ueItem.NetworkType,
			"ril_impl":ueItem.RilImpl,
			"sn":ueItem.SN,
			"baseband":ueItem.Baseband,
			"_save_type":"create",
			"ue_imeis":map[string]interface{}{
				"fieldType":"one2many",
				"modelID":"rt_ue_imei",
				"relatedField":"device_ue_id",
				"list":ueImeis,
			},
		}

		return deviceUe
}

func GetSameSlot(slot *SlotInfo,sameUe map[string]interface{})(map[string]interface{}){
	ueImeiMap,ok:=sameUe["ue_imeis"].(map[string]interface{})
	if !ok {
		return nil
	}

	ueImeiList,ok:=ueImeiMap["list"]
	if !ok || ueImeiList == nil {
		return nil
	}

	for _,ueImei:= range ueImeiList.([]interface{}) {
		if ueImei.(map[string]interface{})["imei"].(string) == slot.IMEI && ueImei.(map[string]interface{})["imsi"].(string) == slot.IMSI {
			return ueImei.(map[string]interface{})
		}
	}
	return nil
}

func UpdateDeviceUeImei(slot *SlotInfo,deviceID string,sameSlot map[string]interface{})(map[string]interface{}){
	return map[string]interface{}{
		"id":sameSlot["id"],
		"version":sameSlot["version"],
		"device_id":deviceID,
		"sn":slot.SlotID,
		"imei":slot.IMEI,
		"imsi":slot.IMSI,
		"operator":slot.Operator,
		"country":slot.Country,
		"short":slot.Short,
		"plmn":slot.PLMN,
		"_save_type":"update",
	}
}

func UpdateDeviceUe(ueItem *UEInfo,deviceID string,sameUe map[string]interface{})(map[string]interface{}){
	ueImeis:=[]map[string]interface{}{}
		for _,slot:= range ueItem.Slot {
			sameSlot:=GetSameSlot(&slot,sameUe)
			var ueImei map[string]interface{}
			if sameSlot == nil {
				ueImei=CreateDeviceUeImei(&slot,deviceID)
			} else {
				ueImei=UpdateDeviceUeImei(&slot,deviceID,sameSlot)
			}
			ueImeis=append(ueImeis,ueImei)
		}

		deviceUe:=map[string]interface{}{
			"id":sameUe["id"],
			"version":sameUe["version"],
			"device_id":deviceID,
			"chip_manufacturer":ueItem.ChipManufacturer,
			"product_manufacturer":ueItem.ProductManufacturer,
			"marketname":ueItem.MarketName,
			//"phone_type":ueItem.PhoneType,
			//"ue_status":ueItem.UeStatus,
			//"network_type":ueItem.NetworkType,
			"ril_impl":ueItem.RilImpl,
			"sn":ueItem.SN,
			"baseband":ueItem.Baseband,
			"_save_type":"update",
			"ue_imeis":map[string]interface{}{
				"fieldType":"one2many",
				"modelID":"rt_ue_imei",
				"relatedField":"device_ue_id",
				"list":ueImeis,
			},
		}

		return deviceUe
}

func GetUpdateDeviceUes(equipmentStatus *EquipmentStatus,deviceItem map[string]interface{})([]map[string]interface{}){
	device_ues:=[]map[string]interface{}{}
	for _,ueItem:= range equipmentStatus.UEDetails {
		sameUe:=GetSameUe(ueItem.SN,deviceItem)
		var deviceUe map[string]interface{}
		if sameUe == nil {
			//create
			deviceUe=CreateDeviceUe(&ueItem,equipmentStatus.RobotId)
		} else {
			deviceUe=UpdateDeviceUe(&ueItem,equipmentStatus.RobotId,sameUe)
		}
		device_ues=append(device_ues,deviceUe)
	}
	return device_ues
}