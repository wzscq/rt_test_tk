import {useSelector,useDispatch} from 'react-redux';
import { useEffect, useState } from 'react';
import {Slider,Button,Space,Select } from 'antd';
import { CaretRightOutlined,PauseOutlined  } from '@ant-design/icons';

import {setCurrentPos} from '../../../redux/dataSlice';
import {createGetTestFileContentMessage} from '../../../utils/normalOperations';


import './index.css';

const options=[
  {
    value: 1000,
    label: '1000',
  },
  {
    value: 2000,
    label: '2000',
  },
  {
    value: 3000,
    label: '3000',
  },
  {
    value: 4000,
    label: '4000',
  },
  {
    value: 5000,
    label: '5000',
  }
];

export default function Header({sendMessageToParent}){
  const dispatch=useDispatch();
  const {origin,item:frameItem}=useSelector(state=>state.frame);
  const device=useSelector(state=>state.data.device.list?state.data.device.list[0]:undefined);
  const currentPos=useSelector(state=>state.data.currentPos);
  const currentRobotInfo=useSelector(state=>state.data.currentRobotInfo);
  const [play,setPlay]=useState(false);
  const [speed,setSpeed]=useState(1000);

  console.log('currentPos',currentPos);

  const onChange=(value)=>{
    dispatch(setCurrentPos(value));
  }

  useEffect(()=>{
    if(device){
      const frameParams={
        frameType:frameItem.frameType,
        frameID:frameItem.params.key,
        origin:origin
      };
      sendMessageToParent(createGetTestFileContentMessage(frameParams,{deviceID:device.device_id,timestamp:device.timestamp,from:parseInt(currentPos)-1,to1:parseInt(currentPos)-1}));
    }
  },[currentPos]);

  useEffect(()=>{
    if(play===true){
      if(currentPos<device.line_count){
        setTimeout(()=>{
          dispatch(setCurrentPos(currentPos+1));
        },speed);
      } else {
        setPlay(false);
      }
    }
  },[play,speed,currentRobotInfo]);

  const convertTimestamptoTime=(unixTimestamp)=>{
    // Convert to milliseconds and
    // then create a new Date object
    let dateObj = new Date(unixTimestamp * 1000);
    let utcString = dateObj.toTimeString();
    
    return utcString.substring(0,8);
}

  const curTime=convertTimestamptoTime(currentRobotInfo?.pcTime);
  //unix time to 

  return (
    <div className='monitor-header'>
      <div className='title'>{'Device: '+device?.device_id+" Start Time: "+device?.start_time}</div>
      <div className='control'>
        <Space>
          <div>{curTime}</div>
          <Select size={'small'}  onChange={(value)=>setSpeed(value)} value={speed} options={options}/>
          <Button size={'small'} type="primary" onClick={()=>setPlay(!play)} icon={play===false?<CaretRightOutlined/>:<PauseOutlined />}></Button>
        </Space>
      </div>
      <div className='slider'>
        <Slider defaultValue={1} dots={false} value={currentPos} max={device.line_count} min={1} onChange={onChange} railStyle={{backgroundColor:"#FFF"}}/>
      </div>
    </div>
  );
}