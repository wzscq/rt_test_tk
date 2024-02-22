import {useSelector} from 'react-redux';
import { useEffect } from 'react';
import useFrame from "../../hook/useFrame";
import Header from "./Header"
import Content from "./Content"

import './index.css';
import PageLoading from "./PageLoading";
import { 
  createQueryDataMessage } from '../../utils/normalOperations';

const queryFields=[
  {field:'id'},
  {field:'device_id'},
  {field:'timestamp'},
  {field:'start_time'},
  {field:'line_count'}
];

export default function Monitor(){
  const sendMessageToParent=useFrame();
  const {origin,item}=useSelector(state=>state.frame);
  const deviceLoaded=useSelector(state=>state.data.deviceLoaded);
  
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
    }
  },[deviceLoaded,item,origin,sendMessageToParent]);

  return (
  <div className="monitor-main">
    {deviceLoaded?
      (<>
        <Header sendMessageToParent={sendMessageToParent}/>
        <Content sendMessageToParent={sendMessageToParent}/>
       </>):
      (<PageLoading/>)}
  </div>);
}