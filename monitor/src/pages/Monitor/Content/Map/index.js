import { useState,useCallback,useEffect } from 'react';
import {useSelector} from 'react-redux';

import Header from './Header';
import Content from './Content';
import {FRAME_MESSAGE_TYPE} from '../../../../utils/constant';

import './index.css';

export default function Map({sendMessageToParent}){
    const {origin,item:frameItem}=useSelector(state=>state.frame);
    const robotMap=useSelector(state=>{
        if(state.data.device?.list&&
           state.data.device?.list.length>0){
            return state.data.device.list[0].robot_map?.list[0];
        }
        return undefined;
    });

    const [fileList,setFileList]=useState([]);

    const getOriginImage=useCallback((files)=>{
        const frameParams={
            frameType:frameItem.frameType,
            frameID:frameItem.params.key,
            dataKey:robotMap.id,
            origin:origin
        };
        const message={
            type:FRAME_MESSAGE_TYPE.GET_IMAGE,
            data:{
                frameParams:frameParams,
                queryParams:{
                    list:files
                }
            }
        }
        sendMessageToParent(message);
    },[sendMessageToParent,robotMap,origin,frameItem]);

    useEffect(()=>{
        const queryResponse=(event)=>{
            const {type,dataKey,data}=event.data;
            if(type===FRAME_MESSAGE_TYPE.QUERY_RESPONSE&&
                dataKey===robotMap.id&&
                data.list&&data.list.length>0){
                //const file=data.list[0];
                const newFileList=robotMap?.file?.list?.map(item=>{
                    const file=data.list.find(element=>element.id===item.id);
                    if(file){
                        return {...item,url:file.url};
                    }
                    return item;
                });
                console.log("newFileList:",newFileList);
                setFileList(newFileList);
            }
        }
        window.addEventListener("message",queryResponse);

        return ()=>{
            window.removeEventListener("message",queryResponse);
        }
    },[robotMap,setFileList]);

    useEffect(()=>{
        if(robotMap&&robotMap.file&&robotMap.file.list){
            //获取图片文件内容
            getOriginImage(robotMap.file.list);
        }
    },[robotMap,getOriginImage]);
    
    return (
        <div className='monitor-map'>
            <Header mapInfo={robotMap}/>
            <Content mapInfo={robotMap} map={fileList[0]}/>
        </div>
    );
}