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
}

func (lfc *LogFileController) update(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("DeviceController runTestCase wrong request")
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

func (lfc *LogFileController) Bind(router *gin.Engine) {
	log.Println("Bind DeviceController")
	router.POST("/logfile/update", lfc.update)
}

func InitLogFileController(conf *common.Config,crvClient *crv.CRVClient,router *gin.Engine){
	dc:=LogFileController{
		CRVClient: crvClient,
		LogFilePath: conf.TestLogFile.Path,
	}

	dc.Bind(router)
}