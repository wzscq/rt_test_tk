import {useSelector} from 'react-redux';
import { useEffect } from 'react';
import useFrame from "../../hook/useFrame";
import Header from "./Header"
import Content from "./Content"

import './index.css';
import PageLoading from "./PageLoading";
import { createGetServerConfMessage } from '../../utils/normalOperations';

export default function Monitor(){
  const sendMessageToParent=useFrame();
  const {origin,item}=useSelector(state=>state.frame);
  const mqttConfLoaded=useSelector(state=>state.mqtt.mqttConfLoaded);

  useEffect(()=>{
    if(mqttConfLoaded===false&&item&&origin){
      const frameParams={
        frameType:item.frameType,
        frameID:item.params.key,
        origin:origin
      };
      sendMessageToParent(createGetServerConfMessage(frameParams));
    }
  },[mqttConfLoaded,item,origin,sendMessageToParent]);

  return (
  <div className="monitor-main">
    {mqttConfLoaded?
      (<><Header/><Content sendMessageToParent={sendMessageToParent}/></>):
      (<PageLoading/>)}
  </div>);
}