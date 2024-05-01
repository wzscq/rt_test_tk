package common

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result interface{} `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type CommonError struct {
	ErrorCode int `json:"errorCode"`
	Params map[string]interface{} `json:"params"`
}

const (
	ResultSuccess = 10000000
	ResultWrongRequest = 10000001
	ResultGetRobotPlatformTokenError = 10200002
	ResultGetRobotPlatformAPIError = 10200003
	ResultMqttClientError = 10200004
	ResultSaveDataError = 10200005
	ResultQueryRequestError = 10200006
	ResultQueryRobotError = 10200007
	ResultNoCommitedTestCaseError = 10200008
	ResultNoTask = 10200009
	ResultNoTaskUe = 10200010
	ResultNoTaskUeTc = 10200011
	ResultConvertPGM2PNGError = 10200012
	ResultFileNotExist = 10200013
	ResultSendTestCaseError = 10200014
	ResultInvokeDeviceAPIError = 10200015
	ResultUpdateLogFileError = 10200016
	ResultDecodeLogFileError = 10200017
	ResultDownloadFileError = 10200018
	ResultTestCaseIsRunning = 10200019
	ResultGetRunningDecodingError = 10200020
	ResultHasRunningDecoding = 10200021
)

var errMsg = map[int]CommonRsp{
	ResultSuccess:CommonRsp{
		ErrorCode:ResultSuccess,
		Message:"操作成功",
		Error:false,
	},
	ResultWrongRequest:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"请求参数错误，请检查参数是否完整，参数格式是否正确",
		Error:true,
	},
	ResultSendTestCaseError:CommonRsp{
		ErrorCode:ResultSendTestCaseError,
		Message:"下发测试用例失败，请与管理员联系处理",
		Error:true,
	},
	ResultGetRobotPlatformTokenError:CommonRsp{
		ErrorCode:ResultGetRobotPlatformTokenError,
		Message:"获取机器人平台授权token失败",
		Error:true,
	},
	ResultGetRobotPlatformAPIError:CommonRsp{
		ErrorCode:ResultGetRobotPlatformAPIError,
		Message:"调用机器人平台API接口发生错误",
		Error:true,
	},
	ResultMqttClientError:CommonRsp{
		ErrorCode:ResultMqttClientError,
		Message:"连接MQTT失败，请与管理员联系处理",
		Error:true,
	},
	ResultSaveDataError:CommonRsp{
		ErrorCode:ResultSaveDataError,
		Message:"保存数据到数据时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultQueryRequestError:CommonRsp{
		ErrorCode:ResultQueryRequestError,
		Message:"下发参数时发送查询参数请求失败，请与管理员联系处理",
		Error:true,
	},
	ResultQueryRobotError:CommonRsp{
		ErrorCode:ResultQueryRobotError,
		Message:"未能查询到对应机器人信息，请与管理员联系处理",
		Error:true,
	},
	ResultNoCommitedTestCaseError:CommonRsp{
		ErrorCode:ResultNoCommitedTestCaseError,
		Message:"未能查询到可发布的测试用例，请与管理员联系处理",
		Error:true,
	},
	ResultNoTask:CommonRsp{
		ErrorCode:ResultNoTask,
		Message:"未能查询到对应的测试任务信息，请与管理员联系处理",
		Error:true,
	},
	ResultNoTaskUe:CommonRsp{
		ErrorCode:ResultNoTaskUe,
		Message:"未获取到测试任务中的UE信息，请检查UE配置是否完整",
		Error:true,
	},
	ResultNoTaskUeTc:CommonRsp{
		ErrorCode:ResultNoTaskUeTc,
		Message:"未获取到测试任务中的UE测试用例信息，请检查测试用例配置是否完整",
		Error:true,
	},
	ResultConvertPGM2PNGError:CommonRsp{
		ErrorCode:ResultConvertPGM2PNGError,
		Message:"PGM格式图片转换为PNG格式图片失败，请检查图片格式是否正确",
		Error:true,
	},
	ResultFileNotExist:CommonRsp{
		ErrorCode:ResultFileNotExist,
		Message:"请求的测试文件不存在",
		Error:true,
	},
	ResultInvokeDeviceAPIError:CommonRsp{
		ErrorCode:ResultInvokeDeviceAPIError,
		Message:"调用设备API接口发生错误",
		Error:true,
	},
	ResultUpdateLogFileError:CommonRsp{
		ErrorCode:ResultUpdateLogFileError,
		Message:"更新日志文件信息失败",
		Error:true,
	},
	ResultDecodeLogFileError:CommonRsp{
		ErrorCode:ResultDecodeLogFileError,
		Message:"解析日志文件失败",
		Error:true,
	},
	ResultDownloadFileError:CommonRsp{
		ErrorCode:ResultDownloadFileError,
		Message:"下载文件失败",
		Error:true,
	},
	ResultTestCaseIsRunning:CommonRsp{
		ErrorCode:ResultTestCaseIsRunning,
		Message:"已经有测试用例正在执行中，请稍后再试",
		Error:true,
	},
	ResultGetRunningDecodingError:CommonRsp{
		ErrorCode:ResultGetRunningDecodingError,
		Message:"获取正在解码的任务数量失败",
		Error:true,
	},
	ResultHasRunningDecoding:CommonRsp{
		ErrorCode:ResultHasRunningDecoding,
		Message:"当前有正在解码的任务，请稍后再试",
		Error:true,
	},
}

func CreateResponse(err *CommonError,result interface{})(*CommonRsp){
	if err==nil {
		commonRsp:=errMsg[ResultSuccess]
		commonRsp.Result=result
		return &commonRsp
	}

	commonRsp:=errMsg[err.ErrorCode]
	commonRsp.Result=result
	commonRsp.Params=err.Params
	return &commonRsp
}

func CreateError(errorCode int,params map[string]interface{})(*CommonError){
	return &CommonError{
		ErrorCode:errorCode,
		Params:params,
	}
}