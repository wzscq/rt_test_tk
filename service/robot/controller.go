package robot

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
	"encoding/base64"
	"github.com/jbuchbinder/gopnm"

	"image/png"
	"bytes"
)

type RobotController struct {
	RobotClient *RobotClient
	CRVClient *crv.CRVClient
	FtpConf *common.FtpConf
}

func (rtc *RobotController)getRobotList(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getRobotList wrong request")
		return
	}	

	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController getRobotList Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if oauthRsp.Result==nil {
		log.Println("RobotController getRobotList Oauth error",oauthRsp)
		params:=map[string]interface{}{
			"message":oauthRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//获取机器人列表
	getRobotListRsp, err:=rtc.RobotClient.GetRobotList(oauthRsp.Result.Token)
	if err!=nil {
		log.Println("RobotController getRobotList GetRobotList error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if getRobotListRsp.Result==nil {
		log.Println("RobotController getRobotList GetRobotList error",getRobotListRsp)
		params:=map[string]interface{}{
			"message":getRobotListRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//机器人信息写入数据库
	rsp:=UpdateRobotList(rtc.CRVClient,getRobotListRsp.Result,header.Token)
	if rsp!=nil {
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

func (rtc *RobotController)getCurrentRobotStatus(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus with error")
		return
  }	

	if rep.SelectedRowKeys ==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getCurrentRobotStatus：request SelectedRowKeys is empty")
		return
	}

	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController getCurrentRobotStatus：request Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if oauthRsp.Result==nil {
		log.Println("RobotController getCurrentRobotStatus：request Oauth error",oauthRsp)
		params:=map[string]interface{}{
			"message":oauthRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//循环获取机器人状态
	for _,robotID:=range *rep.SelectedRowKeys {
		getCurrentRobotStatusRsp, err:=rtc.RobotClient.GetCurrentRobotStatus(oauthRsp.Result.Token,robotID)
		if err!=nil {
			log.Println("RobotController getCurrentRobotStatus：request GetCurrentRobotStatus error",err)
			params:=map[string]interface{}{
				"error":err,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}
	
		if getCurrentRobotStatusRsp.Result==nil {
			log.Println("RobotController getCurrentRobotStatus：request GetCurrentRobotStatus error",getCurrentRobotStatusRsp)
			params:=map[string]interface{}{
				"message":getCurrentRobotStatusRsp.Message,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}

		rsp:=UpdateRobotStatus(rtc.CRVClient,getCurrentRobotStatusRsp.Result,header.Token)
		if rsp!=nil {
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getCurrentRobotStatus success")
}

func (rtc *RobotController)getTestEquipmentStatus(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getTestEquipmentStatus wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getTestEquipmentStatus with error")
		return
  }	

	if rep.SelectedRowKeys ==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController getTestEquipmentStatus：request SelectedRowKeys is empty")
		return
	}

	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController getTestEquipmentStatus：request Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if oauthRsp.Result==nil {
		log.Println("RobotController getTestEquipmentStatus：request Oauth error",oauthRsp)
		params:=map[string]interface{}{
			"message":oauthRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//循环获取机器人状态
	for _,robotID:=range *rep.SelectedRowKeys {
		//获取测试设备信息
		getTestEquipmentStatusRsp, err:=rtc.RobotClient.GetTestEquipmentStatus(oauthRsp.Result.Token,robotID)
		if err!=nil {
			log.Println("RobotController getTestEquipmentStatus：request GetTestEquipmentStatus error",err)
			params:=map[string]interface{}{
				"error":err,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}

		if getTestEquipmentStatusRsp.Result==nil {
			log.Println("RobotController getTestEquipmentStatus：request GetTestEquipmentStatus error",getTestEquipmentStatusRsp)
			params:=map[string]interface{}{
				"message":getTestEquipmentStatusRsp.Message,
			}
			rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return	
		}

		rsp:=UpdateRobotEquipment(rtc.CRVClient,getTestEquipmentStatusRsp.Result,header.Token)
		if rsp!=nil {
			c.IndentedJSON(http.StatusOK, rsp)
			return
		}
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getTestEquipmentStatus success")
}

func (rtc *RobotController)mapUpload(c *gin.Context){
	var header ServerHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload wrong request")
		return
	}	

	var rep UploadMapReq
	if err := c.ShouldBind(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload with error")
		return
  }	

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController mapUpload with error")
		return
	}

	f, _ := file.Open()

	// 解码PGM文件
	//img, format, err := image.Decode(f)
	//log.Println("RobotController mapUpload format",format)
	img, err := pnm.Decode(f)
	if err != nil {
			log.Println("RobotController mapUpload pnm.Decode error",err)
			rsp:=common.CreateResponse(common.CreateError(common.ResultConvertPGM2PNGError,nil),nil)
			c.IndentedJSON(http.StatusOK, rsp)
			return
	}

	content := new(bytes.Buffer)
  // Encode the new image to the buffer in png format.
  encodeErr := png.Encode(content, img)
	if encodeErr != nil {
		log.Println("RobotController mapUpload png.Encode error",encodeErr)
		rsp:=common.CreateResponse(common.CreateError(common.ResultConvertPGM2PNGError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	// 将字节切片转换为Base64编码
	encoded := base64.StdEncoding.EncodeToString(content.Bytes())

	//上传服务器
	rsp:=SaveRobotMap(rtc.CRVClient,&rep,encoded,header.Token)
	if rsp != nil {
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getRobotList success")
}

func (rtc *RobotController)sendTask(c *gin.Context){
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController sendTask wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController sendTask with error")
		return
  }	

	if rep.SelectedRowKeys ==nil || len(*rep.SelectedRowKeys)==0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("RobotController sendTask：request SelectedRowKeys is empty")
		return
	}

	//获取下发内容
	taskID:=(*rep.SelectedRowKeys)[0]
	sendTask,errorCode:=GetTask(rtc.CRVClient,taskID,header.Token,rtc.FtpConf)
	if errorCode != common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	
	//获取token
	oauthRsp, err:=rtc.RobotClient.Oauth()
	if err!=nil {
		log.Println("RobotController sendTask：request Oauth error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformTokenError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	//下发任务配置
	sendTaskRsp, err:=rtc.RobotClient.SendTask(oauthRsp.Result.Token,sendTask)
	if err!=nil {
		log.Println("RobotController sendTask：request SendTask error",err)
		params:=map[string]interface{}{
			"error":err,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	if sendTaskRsp.Success == false {
		log.Println("RobotController sendTask：request SendTask error",sendTaskRsp)
		params:=map[string]interface{}{
			"message":sendTaskRsp.Message,
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRobotPlatformAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return	
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
}

//Bind bind the controller function to url
func (rtc *RobotController) Bind(router *gin.Engine) {
	log.Println("Bind RobotController")
	router.POST("/robot/getRobotList", rtc.getRobotList)
	router.POST("/robot/getCurrentRobotStatus", rtc.getCurrentRobotStatus)
	router.POST("/robot/getTestEquipmentStatus", rtc.getTestEquipmentStatus)
	router.POST("/robot/mapUpload", rtc.mapUpload)
	router.POST("/robot/sendTask", rtc.sendTask)
}