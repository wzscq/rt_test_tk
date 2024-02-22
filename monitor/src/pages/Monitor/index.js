import {useSelector} from 'react-redux';
import { useEffect } from 'react';
import useFrame from "../../hook/useFrame";
import Header from "./Header"
import Content from "./Content"

import './index.css';
import PageLoading from "./PageLoading";
import { createQueryDataMessage,createGetServerConfMessage } from '../../utils/normalOperations';

const queryFields=[
  {field:'id'},
  {field:'host_status'},
  {field:'host_id'},
  {
    field:'robot_map',
    fieldType:'one2many',
    relatedModelID:'rt_robot_map',
    relatedField:'robot_id',
    sorter:[{field:'id',order:'desc'}],
    pagination:{current:1,pageSize:1},
    fields:[
      {field:'id'},
      {field:'robot_id'},
      {field:'picture_id'},
      {field:'picture_name'},
      {field:'building_code'},
      {field:'floor'},
      {
        field:'file',
        fieldType:'file',
      },
    ]
  }
]

export default function Monitor(){
  const sendMessageToParent=useFrame();
  const {origin,item}=useSelector(state=>state.frame);
  const deviceLoaded=useSelector(state=>state.data.deviceLoaded);
  const mqttConfLoaded=useSelector(state=>state.mqtt.mqttConfLoaded);

  useEffect(()=>{
    if(deviceLoaded===false){
    //目前的表单页面仅支持单条数据的编辑和展示
        const dataID=item?.input?.selectedRowKeys[0];
        console.log("dataID:",item);
        if(dataID){
            const frameParams={
                frameType:item.frameType,
                frameID:item.params.key,
                origin:origin
            };
            const queryParams={
                modelID:item.input.modelID,
                filter:{id:dataID},
                fields:queryFields,
                pagination:{current:1,pageSize:1}
            };
            sendMessageToParent(createQueryDataMessage(frameParams,queryParams));
        }
    } else if(mqttConfLoaded===false){
      const frameParams={
        frameType:item.frameType,
        frameID:item.params.key,
        origin:origin
      };
      sendMessageToParent(createGetServerConfMessage(frameParams));
    }
  },[deviceLoaded,mqttConfLoaded,item,origin,sendMessageToParent]);


  return (
  <div className="monitor-main">
    {deviceLoaded&&mqttConfLoaded?
      (<><Header/><Content sendMessageToParent={sendMessageToParent}/></>):
      (<PageLoading/>)}
  </div>);
}