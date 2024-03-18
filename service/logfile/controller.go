package logfile

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
)

type LogFileController struct {
	CRVClient *crv.CRVClient
	LogFilePath string
	DC *DecoderClient
}

func (lfc *LogFileController) update(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController update wrong request")
		return
	}	

	files, err := GetLogFileList(lfc.LogFilePath)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultUpdateLogFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController update with error")
		return
	}

	err = UpdateLogFilesToDB(files, lfc.CRVClient, header.Token)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultUpdateLogFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController update with error")
		return
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("LogFileController update success")
}

func (lfc *LogFileController) parse(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController parse wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.ShouldBind(&rep); err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController parse with error")
		log.Println(err)
		return
  	}

	if rep.SelectedRowKeys == nil || len(*rep.SelectedRowKeys) == 0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController parse with error: SelectedRowKeys is empty")
		return
	}

	err:=DecodeLogFile(rep.SelectedRowKeys, lfc.DC,lfc.CRVClient, header.Token)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultDecodeLogFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController parse with error")
		return
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("LogFileController parse success")
}

func (lfc *LogFileController) Bind(router *gin.Engine) {
	log.Println("Bind DeviceController")
	router.POST("/logfile/update", lfc.update)
	router.POST("/logfile/parse", lfc.parse)
}

func InitLogFileController(conf *common.Config,crvClient *crv.CRVClient,router *gin.Engine){
	dc:=LogFileController{
		CRVClient: crvClient,
		LogFilePath: conf.TestLogFile.Path,
		DC: &DecoderClient{
			URL: conf.TestLogFile.DecoderUrl,
		},
	}

	dc.Bind(router)
}