import { useState,useCallback,useEffect } from 'react';
import {useSelector} from 'react-redux';

import Header from './Header';
import Content from './Content';
import {FRAME_MESSAGE_TYPE} from '../../../../utils/constant';
import { createQueryMapMessage } from '../../../../utils/normalOperations';

import './index.css';

const queryMapFields=[
    {field:'id'},
    {field:'robot_id'},
    {field:'picture_id'},
    {field:'picture_name'},
    {field:'building_code'},
    {field:'floor'},
    {
      field:'file',
      fieldType:'file',
    }
]

export default function Map({sendMessageToParent}){
    const {origin,item:frameItem}=useSelector(state=>state.frame);
    const robotMap=useSelector(state=>state.data.robot_map);
    const robotMapRec=useSelector(state=>{
        console.log('robotMapRec',state.data.robot_map_record);
        if(state.data.robot_map_record?.list&&
            state.data.robot_map_record?.list.length>0){
             return state.data.robot_map_record.list[0];
         }
         return undefined;
        });

    const [fileList,setFileList]=useState([]);

    useEffect(()=>{
        //当robotMap发生变化时，查询robotMap对应的文件信息
        if(robotMap.robot_id&&robotMap.map_id){
            const frameParams={
                frameType:frameItem.frameType,
                frameID:frameItem.params.key,
                origin:origin
            };
            const queryParams={
                modelID:'rt_robot_map',
                filter:{robot_id:robotMap.robot_id,picture_id:robotMap.map_id},
                fields:queryMapFields,
                sorter:[{field:'id',order:'desc'}],
                pagination:{current:1,pageSize:1}
            };
            sendMessageToParent(createQueryMapMessage(frameParams,queryParams));
        }
    },[robotMap]);

    const getOriginImage=useCallback((files)=>{
        const frameParams={
            frameType:frameItem.frameType,
            frameID:frameItem.params.key,
            dataKey:robotMapRec.id,
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
    },[sendMessageToParent,robotMapRec,origin,frameItem]);

    useEffect(()=>{
        const queryResponse=(event)=>{
            const {type,dataKey,data}=event.data;
            if(type===FRAME_MESSAGE_TYPE.QUERY_RESPONSE&&
                dataKey===robotMapRec?.id&&
                data.list&&data.list.length>0){
                //const file=data.list[0];
                const newFileList=robotMapRec?.file?.list?.map(item=>{
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
    },[robotMapRec,setFileList]);

    useEffect(()=>{
        if(robotMapRec&&robotMapRec.file&&robotMapRec.file.list){
            console.log('robotMapRec',robotMapRec);
            //获取图片文件内容
            getOriginImage(robotMapRec.file.list);
        }
    },[robotMapRec,getOriginImage]);
    
    return (
        <div className='monitor-map'>
            <Header mapInfo={robotMapRec} sendMessageToParent={sendMessageToParent}/>
            <Content mapInfo={robotMapRec} map={fileList[0]}/>
        </div>
    );
}