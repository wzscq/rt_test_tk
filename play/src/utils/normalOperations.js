import {OP_TYPE,FRAME_MESSAGE_TYPE,DATA_TYPE} from './constant';

const  GET_FORM_CONF_URL="/definition/getModelFormConf";

const opUpdateModelConf={
    type:OP_TYPE.UPDATE_FRAME_DATA,
    params:{
        dataType:DATA_TYPE.MODEL_CONF
    }
}

const opGetModelConf={
    type:OP_TYPE.REQUEST,
    params:{
        url:GET_FORM_CONF_URL,
        method:"post"
    },
    input:{},
    description:{key:'page.crvformview.getFormConfig',default:'获取模型表单配置信息'}
}

export function createGetFormConfMessage(frameParams,modelID,formID){
    opUpdateModelConf.params={...opUpdateModelConf.params,...frameParams};
    opGetModelConf.input={modelID:modelID,formID:formID};
    opGetModelConf.successOperation=opUpdateModelConf;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opGetModelConf
        }
    };
}

const DATA_QUERY_URL="/data/query";

const opUpdateData={
    type:OP_TYPE.UPDATE_FRAME_DATA,
    params:{
        dataType:DATA_TYPE.QUERY_RESULT
    }
}

const opQueryData={
    type:OP_TYPE.REQUEST,
    params:{
        url:DATA_QUERY_URL,
        method:"post"
    },
    input:{},
    description:{key:'page.crvformview.queryData',default:'查询模型数据'}
}

/**
 * 查询参数如下
 queryParams={modelID,viewID,filter,pagination,sorter,fields}
 */
export function createQueryDataMessage(frameParams,queryParams){
    opUpdateData.params={...opUpdateData.params,...frameParams};
    opQueryData.input=queryParams;
    opQueryData.successOperation=opUpdateData;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opQueryData
        }
    };   
}

const opDownloadFile={
    type:OP_TYPE.DOWNLOAD_FILE,
    params:{
        fileName:"downloadFile",
    },
    input:{},
    description:{key:'page.crvformview.downloadFile',default:'下载文件'}
}

export function createDownloadFileMessage(file,fileName){
    opDownloadFile.input=file;
    opDownloadFile.params.fileName=fileName;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opDownloadFile
        }
    };   
}


const opUpdateFileContent={
    type:OP_TYPE.UPDATE_FRAME_DATA,
    params:{
        dataType:DATA_TYPE.TEST_FILE_CONENT
    }
}

const REDIRECT_URL='/redirect';
const RT_SERVICE_GET_FILE_CONTENT='rtservice/testfile/GetContent';
const RT_SERVICE_GET_FILE_POINTS='rtservice/testfile/GetPoints';

const opGetFileContent={
    type:OP_TYPE.REQUEST,
    params:{
        url:REDIRECT_URL,
        method:"post"
    },
    input:{},
    description:{key:'page.rtservice.play.getfilecontent',default:'获取测试文件数据'},
    queenable:true
}

export function createGetTestFileContentMessage(frameParams,device){
    opUpdateFileContent.params={...opUpdateFileContent.params,...frameParams};
    opGetFileContent.input={...device,to:RT_SERVICE_GET_FILE_CONTENT};
    opGetFileContent.successOperation=opUpdateFileContent;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opGetFileContent
        }
    };   
}

const opUpdateFilePoints={
    type:OP_TYPE.UPDATE_FRAME_DATA,
    params:{
        dataType:DATA_TYPE.TEST_FILE_POINTS
    }
}

const opGetFilePoints={
    type:OP_TYPE.REQUEST,
    params:{
        url:REDIRECT_URL,
        method:"post"
    },
    input:{},
    description:{key:'page.rtservice.play.getfilecontent',default:'获取测试文件数据'},
    queenable:true
}

export function createGetTestFilePointsMessage(frameParams,device){
    opUpdateFilePoints.params={...opUpdateFilePoints.params,...frameParams};
    opGetFilePoints.input={...device,to:RT_SERVICE_GET_FILE_POINTS};
    opGetFilePoints.successOperation=opUpdateFilePoints;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opGetFilePoints
        }
    };
}

const opUpdateMap={
    type:OP_TYPE.UPDATE_FRAME_DATA,
    params:{
        dataType:DATA_TYPE.ROBOT_MAP_RECORD
    }
}

export function createQueryMapMessage(frameParams,queryParams){
    opUpdateMap.params={...opUpdateMap.params,...frameParams};
    opQueryData.input=queryParams;
    opQueryData.successOperation=opUpdateMap;
    opQueryData.queenable=true;
    return {
        type:FRAME_MESSAGE_TYPE.DO_OPERATION,
        data:{
            operationItem:opQueryData
        }
    };   
}