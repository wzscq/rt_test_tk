import {useSelector} from 'react-redux';
import { SplitPane } from "react-collapse-pane";
import PropertyGrid from './PropertyGrid';
import UEContent from './UEContent';
import Map from './Map';

import './index.css';
import { useMemo } from 'react';

export default function Content({sendMessageToParent}){
  //const currentRobotInfo=useSelector(state=>state.data.currentRobotInfo);
  const currentUes=useSelector(state=>state.data.currentUes);
  
  const ueControls=useMemo(()=>{
    if(Object.keys(currentUes).length>0){
      const imsiList=Object.keys(currentUes).sort((a,b)=>a>b);
      console.log("imsiList",imsiList);
      return imsiList.map(imsi => {
        return (<UEContent key={imsi} imsi={imsi} />)
      });
    }
    return ([<div></div>]);
  },[currentUes]);

  return (
    <div className='monitor-content'>
      <SplitPane dir='ltr'initialSizes={[40,60]} split="vertical" collapse={false}>
        <Map sendMessageToParent={sendMessageToParent}/>    
        <div className='monitor-content-right'>
          <SplitPane dir='rtl' split="horizontal" collapse={false}>
            {ueControls}
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  );

  /*return (
    <div className='monitor-content'>
      <SplitPane dir='ltr'initialSizes={[15,25,60]} split="vertical" collapse={false}>
        <div className='monitor-content-left'>
          <PropertyGrid obj={currentRobotInfo} title="robot info"/>
        </div>
        <div className='monitor-content-center'>
        <SplitPane dir='rtl'initialSizes={[60,40]} split="horizontal" collapse={false}>
          <Map sendMessageToParent={sendMessageToParent}/>    
          <div></div>
        </SplitPane>
        </div>
        <div className='monitor-content-right'>
          <SplitPane dir='rtl' split="horizontal" collapse={false}>
            {ueControls}
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  );*/

}