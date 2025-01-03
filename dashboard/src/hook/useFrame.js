import { useEffect,useCallback } from 'react';
import {useSelector,useDispatch} from 'react-redux';

import {setParam} from '../redux/frameSlice';
import {setImsi,setDial,setAttachStatus,setRATStatus} from '../redux/dataSlice';
import {setLocale} from '../redux/i18nSlice';
import {setServerConf} from '../redux/mqttSlice';

import {
    FRAME_MESSAGE_TYPE,
    DATA_TYPE
} from '../utils/constant';

const getParentOrigin=()=>{
    const a = document.createElement("a");
    a.href=document.referrer;
    return a.origin;
}

export default function useFrame(){
    const dispatch=useDispatch();
    const {origin}=useSelector(state=>state.frame);
    //const {forms} = useSelector(state=>state.definition);

    const sendMessageToParent=useCallback((message)=>{
        if(origin){
            window.parent.postMessage(message,origin);
        } else {
            console.log("the origin of parent is null,can not send message to parent.");
        }
    },[origin]);
        
    //这里在主框架窗口中挂载事件监听函数，负责和子窗口之间的操作交互
    const receiveMessageFromMainFrame=useCallback((event)=>{
        console.log("crv_form receiveMessageFromMainFrame:",event);
        const {type,dataType,data}=event.data;
        if(type===FRAME_MESSAGE_TYPE.INIT){
            dispatch(setParam({origin:event.origin,item:event.data.data}));
            if(event.data.i18n){
                dispatch(setLocale(event.data.i18n));
            }
        } else if (type===FRAME_MESSAGE_TYPE.UPDATE_DATA){
            if(dataType===DATA_TYPE.MODEL_CONF){
                //dispatch(setDefinition(data));
            } else if (dataType===DATA_TYPE.QUERY_RESULT){
                //dispatch(setDevice(data));
            } else if (dataType===DATA_TYPE.SERVER_CONF){
                dispatch(setServerConf(data));
            } else if (dataType===DATA_TYPE.IMSI){
                dispatch(setImsi(data));
            } else if (dataType===DATA_TYPE.DIAL){
                dispatch(setDial(data));
            } else if (dataType===DATA_TYPE.ATTACH){
                dispatch(setAttachStatus(data));
            } else if (dataType===DATA_TYPE.ATTACH_QUERY){
                dispatch(setAttachStatus(data));
            } else if (dataType===DATA_TYPE.DETACH){
                dispatch(setAttachStatus(data));
            } else if (dataType===DATA_TYPE.RAT){
                dispatch(setRATStatus(data));
            }else {
                console.log("update data with wrong data type:",dataType);
            }
        } else if (type===FRAME_MESSAGE_TYPE.RELOAD_DATA){
            //dispatch(refreshData());
        } else if (type===FRAME_MESSAGE_TYPE.UPDATE_LOCALE){
            //console.log("UPDATE_LOCALE",event.data)
            //dispatch(setLocale(event.data.i18n));
        }
    },[dispatch]);
        
    useEffect(()=>{
        window.addEventListener("message",receiveMessageFromMainFrame);
        return ()=>{
            window.removeEventListener("message",receiveMessageFromMainFrame);
        }
    },[receiveMessageFromMainFrame]);

    useEffect(()=>{
        if(origin===null){
            setTimeout(()=>{
                console.log('postMessage to parent init');
                window.parent.postMessage({type:FRAME_MESSAGE_TYPE.INIT},getParentOrigin());
            },200);
        }
    },[origin]);

    return sendMessageToParent;
}