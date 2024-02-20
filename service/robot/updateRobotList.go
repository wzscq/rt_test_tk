package robot

import (
	"log"
	"rt_test_service/common"
	"rt_test_service/crv"
	"strings"
)

type ProjectRobot struct {
	ID string `json:"id"`
	Version string `json:"version"`
	Name string `json:"name"`
	ProjectID string `json:"rt_project_id"`
}	

var MODELID_ROBOT="rt_project_robot"
var MODELID_ROBOT_EVENT="rt_robot_event"

var	queryRobotFields=[]map[string]interface{}{
	{"field": "id"},
	{"field": "version"},
}

func GetRobotList(crvClient *crv.CRVClient,token string,filter map[string]interface{})([]ProjectRobot){
	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT,
		Fields:&queryRobotFields,
		Filter:&filter,
	}

	rsp,commonErr:=crvClient.Query(&commonRep,token)
	if commonErr!=common.ResultSuccess {
		return nil
	}

	if rsp.Error == true {
		log.Println("GetRobotList error:",rsp.ErrorCode,rsp.Message)
		return nil
	}

	resLst,ok:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	if !ok {
		log.Println("GetRobotList error: no list in rsp.")
		return nil
	}

	robots:=[]ProjectRobot{}

	for _,res:=range resLst {
		resMap,ok:=res.(map[string]interface{})
		if !ok {
			log.Println("GetRobotList error: can not convert item to map.")
			return nil
		}

		robotItem:=ProjectRobot{
			ID:resMap["id"].(string),
			Version:resMap["version"].(string),
		}

		robots=append(robots,robotItem)
	}

	return robots
}

func UpdateRobotList(crvClient *crv.CRVClient,robotList []RobotInfo,token string)(*common.CommonRsp){
	//获取已经存在的机器人列表
	currentRobotList:=GetRobotList(crvClient,token,nil)

	saveRobotList:=[]map[string]interface{}{}
	//比较机器人列表，找出需要新增的机器人和需要删除的机器人
	//找出需要更新和删除的机器人
	for _,currentRobot:=range currentRobotList {
		//在新的机器人列表中找到当前机器人
		var robotInfo *RobotInfo 
		robotInfo=nil
		for _,robot:=range robotList {
			if currentRobot.ID==robot.RobotId {
				robotInfo=&robot
				break
			}
		}

		if robotInfo != nil {
			//将用逗号分割的经纬度字符串拆分成数组
			llarr:=strings.Split(robotInfo.LongitudeLatitude,",")

			isOnline:="0"
			if robotInfo.IsOnlineStatus==true {
				isOnline="1"
			}

			saveRobot:=map[string]interface{}{
				"id":currentRobot.ID,
				"version":currentRobot.Version,
				"name":robotInfo.RobotName,
				"longitude":llarr[0],
				"latitude":llarr[1],
				"position":robotInfo.Position,
				"electricity":robotInfo.Electricity,
				"online":robotInfo.Online,
				"is_online":isOnline,
				"_save_type":"update",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		} else {
			//需要删除
			saveRobot:=map[string]interface{}{
				"id":currentRobot.ID,
				"version":currentRobot.Version,
				"_save_type":"delete",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		}
	}
	//找出新增的机器人
	for _,robot:=range robotList {
		var currentRobot *ProjectRobot 
		currentRobot=nil
		for _,currentRobotItem:=range currentRobotList {
			if currentRobotItem.ID==robot.RobotId {
				currentRobot=&currentRobotItem
				break
			}
		}
		if currentRobot == nil {
			//将用逗号分割的经纬度字符串拆分成数组
			llarr:=strings.Split(robot.LongitudeLatitude,",")

			isOnline:="0"
			if robot.IsOnlineStatus==true {
				isOnline="1"
			}

			//需要新增
			saveRobot:=map[string]interface{}{
				"id":robot.RobotId,
				"name":robot.RobotName,
				"longitude":llarr[0],
				"latitude":llarr[1],
				"position":robot.Position,
				"electricity":robot.Electricity,
				"online":robot.Online,
				"is_online":isOnline,
				"_save_type":"create",
			}
			saveRobotList=append(saveRobotList,saveRobot)
		}
	}

	//保存机器人列表
	commonRep:=crv.CommonReq{
		ModelID:MODELID_ROBOT,
		List:&saveRobotList,
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

	//保存机器人事件记录
	saveRobotEventList:=[]map[string]interface{}{}
	for _,robot:=range robotList {
		//将用逗号分割的经纬度字符串拆分成数组
		llarr:=strings.Split(robot.LongitudeLatitude,",")

		isOnline:="0"
		if robot.IsOnlineStatus==true {
			isOnline="1"
		}

		//需要新增
		saveRobot:=map[string]interface{}{
			"robot_id":robot.RobotId,
			"robot_name":robot.RobotName,
			"longitude":llarr[0],
			"latitude":llarr[1],
			"position":robot.Position,
			"electricity":robot.Electricity,
			"online":robot.Online,
			"is_online":isOnline,
			"_save_type":"create",
		}
		saveRobotEventList=append(saveRobotEventList,saveRobot)
	}

	//保存机器人列表
	commonRep=crv.CommonReq{
		ModelID:MODELID_ROBOT_EVENT,
		List:&saveRobotEventList,
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

