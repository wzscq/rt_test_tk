import {useSelector} from 'react-redux';

import Header from "./Header";
import Status from './Status';
import useFrame from "../../hook/useFrame";

export default function Dashboard() {
  const sendMessageToParent=useFrame();
  const frame=useSelector(state=>state.frame);


  if(frame.origin){
    return (
      <div>
        <Header sendMessageToParent={sendMessageToParent} frame={frame}/>
        <Status sendMessageToParent={sendMessageToParent} frame={frame}/>
      </div>
    );
  }

  return (
    <div>loading {JSON.stringify(frame)}</div>
  );
}