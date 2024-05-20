package logfile

import (
	"log"
	"github.com/gin-gonic/gin"
	"rt_test_service/common"
	"rt_test_service/crv"
	"net/http"
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

type LogFileController struct {
	CRVClient *crv.CRVClient
	LogFilePath string
	DecodedPath string
	LogFileFromDecoderUrl string
	DC *DecoderClient
	ExpandTimeRange time.Duration
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

	//清空数据库
	DeleteAllLogFiles(lfc.CRVClient,header.Token)

	err = UpdateLogFilesToDB(files, lfc.CRVClient, header.Token, lfc.ExpandTimeRange)
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

func (lfc *LogFileController) updateFromDecoder(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController update wrong request")
		return
	}	

	files, err := GetLogFileListFromDecode(lfc.LogFileFromDecoderUrl)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultUpdateLogFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController update with error")
		return
	}

	//清空数据库
	DeleteAllLogFiles(lfc.CRVClient,header.Token)

	err = UpdateLogFilesToDB(files, lfc.CRVClient, header.Token, lfc.ExpandTimeRange)
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

func (lfc *LogFileController) decode(c *gin.Context) {
	log.Println("LogFileController decode start ...")
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.ShouldBind(&rep); err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode with error")
		log.Println(err)
		return
  	}

	if rep.SelectedRowKeys == nil || len(*rep.SelectedRowKeys) == 0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode with error: SelectedRowKeys is empty")
		return
	}

	//判断当前是否存在正在解码的任务
	runningCount,err:=GetDecodingTaskCount(lfc.CRVClient, header.Token)
	if err!=nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultGetRunningDecodingError,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode with error: GetDecodingTaskCount error")
		return
	}

	if runningCount>0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultHasRunningDecoding,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode with error: there are decoding tasks running")
		return
	}

	err=DecodeLogFile(rep.SelectedRowKeys, lfc.DC,lfc.CRVClient, header.Token)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultDecodeLogFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController decode with error")
		return
	}

	rsp:=common.CreateResponse(nil,nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("LogFileController decode success")
}

func (lfc *LogFileController) download(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.ShouldBind(&rep); err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download with error")
		log.Println(err)
		return
  }

	if rep.SelectedRowKeys == nil || len(*rep.SelectedRowKeys) == 0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download with error: SelectedRowKeys is empty")
		return
	}

	//id:=(*rep.SelectedRowKeys)[0]
	//string to int64 
	id,err:=strconv.ParseInt((*rep.SelectedRowKeys)[0],10,64)

	file,err:=GetDeodeRecordFromDB(id, lfc.CRVClient, header.Token)
	if err != nil {
		params:=map[string]interface{}{
			"error":err.Error(),
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultDownloadFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download with error")
		return
	}

	if file==nil{
		params:=map[string]interface{}{
			"error":"no file",
		}
		rsp:=common.CreateResponse(common.CreateError(common.ResultDownloadFileError,params),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download with error")
		return
	}

	decodedFileName,_:=file["decoded_file"].(string)
	//替换掉文件固定的前缀
	log.Println("decodedFileName:",decodedFileName)
	decodedFileName = filepath.Base(decodedFileName)
	log.Println("decodedFileName without path:",decodedFileName)

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", decodedFileName))

	decodedFileName = filepath.Join(lfc.DecodedPath,decodedFileName)
	log.Println("decodedFileName with path:",decodedFileName)

	c.File(decodedFileName)
}

func (lfc *LogFileController) Bind(router *gin.Engine) {
	log.Println("Bind DeviceController")
	router.POST("/logfile/update", lfc.update)
	router.POST("/logfile/decode", lfc.decode)
	router.POST("/logfile/download", lfc.download)
	router.POST("/logfile/updateFromDecoder", lfc.updateFromDecoder)
}

func InitLogFileController(conf *common.Config,crvClient *crv.CRVClient,router *gin.Engine){
	duration, _ := time.ParseDuration(conf.TestLogFile.ExpandTimeRange)
	dc:=LogFileController{
		CRVClient: crvClient,
		LogFilePath: conf.TestLogFile.Path,
		DecodedPath: conf.TestLogFile.DecodedPath,
		ExpandTimeRange: duration,
		DC: &DecoderClient{
			URL: conf.TestLogFile.DecoderUrl,
			GetLogFileUrl: conf.TestLogFile.GetLogFileUrl,
		},
		LogFileFromDecoderUrl: conf.TestLogFile.LogFileFromDecoderUrl,
	}

	dc.Bind(router)
}