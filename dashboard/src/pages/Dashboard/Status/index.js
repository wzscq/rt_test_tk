import {useSelector} from 'react-redux';

import Imsi from './Imsi';
import Dial from './Dial';
import DeviceControl from './DeviceControl';

export default function Status({sendMessageToParent,frame}) {


  return (
    <div style={{padding:'10px'}}>
      <Imsi sendMessageToParent={sendMessageToParent} frame={frame}/>
      <div style={{height:10}}></div>
      <Dial sendMessageToParent={sendMessageToParent} frame={frame}/>
      <div style={{height:10}}></div>
      <DeviceControl sendMessageToParent={sendMessageToParent} frame={frame}/>
    </div>
  );
}