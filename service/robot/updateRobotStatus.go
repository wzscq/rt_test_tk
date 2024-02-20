package robot

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
	"strings"
)

var MODELID_ROBOT_STATUS="rt_robot_status"

func UpdateRobotStatus(crvClient *crv.CRVClient,robotStatus *RobotStatus,token string)(*common.CommonRsp){
	//获取已经存在的机器人列表
	filter:=map[string]interface{}{
		"id":robotStatus.RobotId,
	}
	currentRobotList:=GetRobotList(crvClient,token,filter)
	if currentRobotList==nil || len(currentRobotList)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultQueryRobotError,nil),nil)
		return rsp
	}

	//更新机器人信息
	robot:=currentRobotList[0]
	electricity:= strings.Replace(robotStatus.Electricity,"%","",-1)

	saveRobot:=map[string]interface{}{
		"id":robot.ID,
		"version":robot.Version,
		"map_id":robotStatus.MapId,
		"current_task":robotStatus.CurrentTask,
		"robot_status":robotStatus.RobotStatus,
		"pixel_x":robotStatus.PixelX,
		"pixel_y":robotStatus.PixelY,
		"pixel_theta":robotStatus.PixelTheta,
		"exception":robotStatus.Exception,
		"map_name":robotStatus.MapName,
		"state_code":robotStatus.StateCode,
		"task_id":robotStatus.TaskId,
		"task_name":robotStatus.TaskName,
		"electricity":electricity,
		"_save_type":"update",
	}

	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT,
		List:&[]map[string]interface{}{saveRobot},
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

	//
	saveRobotStatus:=map[string]interface{}{
		"robot_id":robotStatus.RobotId,
		"robot_name":robotStatus.RobotName,
		"map_id":robotStatus.MapId,
		"current_task":robotStatus.CurrentTask,
		"robot_status":robotStatus.RobotStatus,
		"pixel_x":robotStatus.PixelX,
		"pixel_y":robotStatus.PixelY,
		"pixel_theta":robotStatus.PixelTheta,
		"exception":robotStatus.Exception,
		"map_name":robotStatus.MapName,
		"state_code":robotStatus.StateCode,
		"task_id":robotStatus.TaskId,
		"task_name":robotStatus.TaskName,
		"electricity":electricity,
		"_save_type":"create",
	}

	commonRep=crv.CommonReq{
		ModelID:MODELID_ROBOT_STATUS,
		List:&[]map[string]interface{}{saveRobotStatus},
	}

	rsp,commonErr=crvClient.Save(&commonRep,token)
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