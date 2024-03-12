import { SplitPane } from "react-collapse-pane";
import UEContent from './UEContent';
import Map from './Map';

import './index.css';

export default function Content({sendMessageToParent}){
  return (
    <div className='monitor-content'>
      <SplitPane dir='ltr' initialSizes={[40,60]} split="vertical" collapse={false}>
        <Map sendMessageToParent={sendMessageToParent}/>    
        <div className='monitor-content-right'>
          <UEContent/>
        </div>
      </SplitPane>
    </div>
  );
}