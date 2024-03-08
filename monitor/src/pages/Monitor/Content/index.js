import {useSelector} from 'react-redux';
import { SplitPane } from "react-collapse-pane";
import UEContent from './UEContent';
import Map from './Map';

import './index.css';
import { useMemo } from 'react';

export default function Content({sendMessageToParent}){
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
      <SplitPane dir='ltr' initialSizes={[40,60]} split="vertical" collapse={false}>
        <Map sendMessageToParent={sendMessageToParent}/>    
        <div className='monitor-content-right'>
          <SplitPane dir='rtl' split="horizontal" collapse={false}>
            {ueControls}
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  );
}