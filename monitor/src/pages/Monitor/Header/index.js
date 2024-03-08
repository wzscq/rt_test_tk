import {useSelector,useDispatch} from 'react-redux';
import mqtt from 'mqtt';
import { useEffect, useState } from 'react';

import {addDataItem} from '../../../redux/dataSlice';

import './index.css';

var g_MQTTClient=null;

export default function Header(){
  const dispatch=useDispatch();
  const {mqttConf,mqttStatus}=useSelector(state=>state.mqtt);
  //const [mqttStatus,setMqttStatus]=useState('disconnected');

  useEffect(()=>{
    const connectMqtt=()=>{
      console.log("connectMqtt ... ");
      if(g_MQTTClient!==null){
          g_MQTTClient.end();
          g_MQTTClient=null;
      }
  
      const server='ws://'+mqttConf.broker+':'+mqttConf.wsPort;
      const options={
          username:mqttConf.user,
          password:mqttConf.password,
          keepalive:3600,
          reconnectPeriod:60
      }
      console.log("connect to mqtt server ... "+server+" with options:",options);
      g_MQTTClient  = mqtt.connect(server,options);
      g_MQTTClient.on('connect', () => {
          setMqttStatus("connected to mqtt server "+server+".");
          const topic=mqttConf.uploadMeasurementMetrics;
          g_MQTTClient.subscribe(topic, (err) => {
              if(!err){
                  setMqttStatus("subscribe topics success.");
                  console.log("topic:",topic);
              } else {
                  setMqttStatus("subscribe topics error :"+err.toString());
              }
          });
      });
      g_MQTTClient.on('message', (topic, payload, packet) => {
          console.log("receiconsolleve message topic :"+topic+" content :"+payload.toString());
          dispatch(addDataItem(JSON.parse(payload.toString())));
      });
      g_MQTTClient.on('close', () => {
        setMqttStatus("mqtt client is closed.");
      });
    }

    //connectMqtt();
  },[dispatch,mqttConf]);

  return (
    <div className='monitor-header'>
      {' MQTT: '+mqttStatus}
    </div>
  )
}