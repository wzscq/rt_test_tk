package device

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
)

type DeviceController struct {
	CRVClient *crv.CRVClient
	MqttConf *common.MqttConf
	FtpConf *common.FtpConf
}

func (dc *DeviceController)getServerConf(c *gin.Context){
	res:=map[string]interface{}{
		"mqtt":dc.MqttConf,
		"ftp":dc.FtpConf,
	}
	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("RobotController getServerConf")
}

func (dc *DeviceController)getTestCase(c *gin.Context){
	var header ServerHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController getTestCase wrong request")
		return
	}	

	var rep DeviceReq
	if err := c.ShouldBind(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController getTestCase with error")
		return
  }

	//查询可下发的测试用例
	rsp:=GetCommitedTestCase(dc.CRVClient,header.Token)
	if rsp.Error {
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController getTestCase with error")
		return
	}

	//更新机器人测试用例
	tcList:=rsp.Result.(map[string]interface{})["list"].([]interface{})
	rsp=UpdateRobotTestCase(dc.CRVClient,header.Token,rep.RobotID,tcList)

	if rsp != nil {
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController getTestCase with error")
		return
	}

	//生成下发数据结构
	res:=GetTestCaseForDevice(tcList)

	rsp=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("DeviceController getTestCase success")
}

func (dc *DeviceController) Bind(router *gin.Engine) {
	log.Println("Bind DeviceController")
	router.POST("/device/getServerConf", dc.getServerConf)
	router.POST("/device/getTestCase", dc.getTestCase)
}