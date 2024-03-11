package device

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"rt_test_service/mqtt"
	"net/http"
)

type DeviceController struct {
	CRVClient *crv.CRVClient
	MqttConf *common.MqttConf
	FtpConf *common.FtpConf
	MQTTClient *mqtt.MQTTClient
	MapConf *common.MapConf
	DeviceClient *DeviceClient
}

func (dc *DeviceController)getServerConf(c *gin.Context){
	res:=map[string]interface{}{
		"mqtt":dc.MqttConf,
		"ftp":dc.FtpConf,
		"map":dc.MapConf,
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

func (dc *DeviceController)runTestCase(c *gin.Context){
	var header ServerHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.ShouldBind(&rep); err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase with error")
		log.Println(err)
		return
  }

	if rep.SelectedRowKeys == nil || len(*rep.SelectedRowKeys) == 0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase with error: SelectedRowKeys is empty")
		return
	}

	//查询测试用例
	tc:=GetTestCase((*rep.SelectedRowKeys)[0],header.Token,dc.CRVClient)
	if tc==nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultNoCommitedTestCaseError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase with error: query test case failed")
		return
	}
	
	//生成下发数据结构
	err:=SendTestCase(tc,dc.MQTTClient,dc.MqttConf.SendTestCaseTopic)
	if err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultSendTestCaseError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase with error: send test case failed")
		return
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (dc *DeviceController)getImsi(c *gin.Context){
	log.Println("DeviceController getImsi")
	res,err:=dc.DeviceClient.GetImsi()
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultInvokeDeviceAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController getImsi with error")
		return
	}

	log.Println(res)

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("DeviceController getImsi success")
}

func (dc *DeviceController)dialQuery(c *gin.Context){
	res,err:=dc.DeviceClient.GetDialQuery()
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultInvokeDeviceAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController dialQuery with error")
		return
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("DeviceController dialQuery success")
}

func (dc *DeviceController)dialTrigger(c *gin.Context){
	res,err:=dc.DeviceClient.DialTrigger()
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultInvokeDeviceAPIError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController dialTrigger with error")
		return
	}

	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),res)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("DeviceController dialQuery success")
}

func (dc *DeviceController) Bind(router *gin.Engine) {
	log.Println("Bind DeviceController")
	router.POST("/device/getServerConf", dc.getServerConf)
	router.POST("/device/getTestCase", dc.getTestCase)
	router.POST("/device/runTestCase", dc.runTestCase)
	router.POST("/device/getImsi", dc.getImsi)
	router.POST("/device/dialQuery", dc.dialQuery)
	router.POST("/device/dialTrigger", dc.dialTrigger)
}