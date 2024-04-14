package testfile

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rt_test_service/common"
	"rt_test_service/crv"
	"fmt"
	//"strconv"
)

type GetContentReq struct {
	DeviceID  string `json:"deviceID"`
	TimeStamp string `json:"timestamp"`
	From      int64  `json:"from"`
	To        int64  `json:"to1"`
}

type IndicatorLegendItem struct {
	ID    string `json:"id"`
	SN    string `json:"sn"`
	Start string `json:"start"`
	End   string `json:"end"`
	RGB   string `json:"rgb"`
}

type IndicatorLegend struct {
	ModelID string                `json:"modelID"`
	Total   int                   `json:"total"`
	List    []IndicatorLegendItem `json:"list"`
}

type Indicator struct {
	ExtractPath string          `json:"extract_path"`
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Legend      IndicatorLegend `json:"legend"`
}

type GetPointsReq struct {
	DeviceID  string    `json:"deviceID"`
	TimeStamp string    `json:"timestamp"`
	Indicator Indicator `json:"indicator"`
}

type TestFileController struct {
	OutPath string
	CRVClient *crv.CRVClient
}

func (tfc *TestFileController) GetContent(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent wrong request")
		return
	}

	var rep GetContentReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent with error")
		return
	}

	//获取文件内容
	//timestamp,_:=strconv.ParseInt(rep.TimeStamp,10,64)
	tf := GetTestFile(tfc.OutPath, rep.DeviceID, rep.TimeStamp)
	if tf == nil {
		rsp := common.CreateResponse(common.CreateError(common.ResultFileNotExist, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetContent file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	//获取文件内容
	//from,_:=strconv.ParseInt(rep.From,10,64)
	//to,_:=strconv.ParseInt(rep.To,10,64)
	content := tf.GetContent(rep.From, rep.To)
	rsp := common.CreateResponse(common.CreateError(common.ResultSuccess, nil), content)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (tfc *TestFileController) GetPoints(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints wrong request")
		return
	}

	var rep GetPointsReq
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints with error")
		return
	}

	//获取文件内容
	//timestamp,_:=strconv.ParseInt(rep.TimeStamp,10,64)
	tf := GetTestFile(tfc.OutPath, rep.DeviceID, rep.TimeStamp)
	if tf == nil {
		rsp := common.CreateResponse(common.CreateError(common.ResultFileNotExist, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController GetPoints file not exist.")
		return
	}
	defer tf.CloseReadOnly()

	//获取文件内容
	points := tf.GetPoints(rep.Indicator)
	rsp := common.CreateResponse(common.CreateError(common.ResultSuccess, nil), points)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (tfc *TestFileController) download(c *gin.Context) {
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController download wrong request")
		return
	}	

	var rep crv.CommonReq
	if err := c.ShouldBind(&rep); err != nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController download with error")
		log.Println(err)
		return
  }

	if rep.SelectedRowKeys == nil || len(*rep.SelectedRowKeys) == 0 {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("TestFileController download with error: SelectedRowKeys is empty")
		return
	}

	//id:=(*rep.SelectedRowKeys)[0]
	//string to int64 
	id:=(*rep.SelectedRowKeys)[0]

	file,errorCode:=GetTestFileFromDB(id, header.Token,tfc.CRVClient)
	if errorCode != common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("LogFileController download with error")
		return
	}

	

	fileName:=file.DeviceID+"_"+file.TimeStamp+".zip"
	//替换掉文件固定的前缀
	log.Println("fileName:",fileName)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	fileName = tfc.OutPath +"/"+ fileName
	log.Println("decodedFileName with path:",fileName)

	c.File(fileName)
}

//Bind bind the controller function to url
func (tfc *TestFileController) Bind(router *gin.Engine) {
	log.Println("Bind CacheFileController")
	router.POST("/testfile/GetContent", tfc.GetContent)
	router.POST("/testfile/GetPoints", tfc.GetPoints)
	router.POST("/testfile/download", tfc.download)
}
