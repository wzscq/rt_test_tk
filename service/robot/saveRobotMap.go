package robot

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
)

var MODELID_ROBOT_MAP="rt_robot_map"

func SaveRobotMap(crvClient *crv.CRVClient,rep *UploadMapReq,fileBase64 string,token string)(*common.CommonRsp){
	saveRobotMap:=map[string]interface{}{
		"robot_id":rep.RobotId,
		"picture_id":rep.PictureId,
		"picture_name":rep.PictureName,
		"building_code":rep.BuildingCode,
		"floor":rep.Floor,
		"file":map[string]interface{}{
			"fieldType":"file",
			"list":[]map[string]interface{}{
				map[string]interface{}{
					"contentBase64":fileBase64,
					"name":rep.PictureName+".png",
					"_save_type":"create",
				},
			},
		},
		"_save_type":"create",
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT_MAP,
		List:&[]map[string]interface{}{saveRobotMap},
	}

	rsp,commonErr:=crvClient.Save(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(commonErr,nil),nil)
		return rsp
	}

	if rsp.Error == true {
		log.Println("UpdateRobotList error:",rsp.ErrorCode,rsp.Message)
		return rsp
	}

	return nil
}